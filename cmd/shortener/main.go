package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

type URL struct {
	urls map[string]string
}

var urlStorage = URL{
	urls: make(map[string]string),
}

func generateShortURL() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, postHandler)
	mux.HandleFunc(`/{id}`, getHandler)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Этот URL принимает только POST запросы", http.StatusBadRequest)
		return
	}
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/plain") {
		http.Error(w, "Неверный Content-Type. Ожидается text/plain", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	originalURL := string(body)
	if originalURL == "" {
		http.Error(w, "URL не может быть пустым", http.StatusBadRequest)
		return
	}
	var shortKey string
	for k, v := range urlStorage.urls {
		if v == originalURL {
			shortKey = k
			break
		}
	}
	if shortKey == "" {
		shortKey = generateShortURL()
		urlStorage.urls[shortKey] = originalURL
	}
	shortURL := fmt.Sprintf("http://localhost:8080/%s", shortKey)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Этот URL принимает только GET запросы", http.StatusBadRequest)
		return
	}
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	if shortURL == "" {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	originalURL, exists := urlStorage.urls[shortURL]
	if !exists {
		http.Error(w, "Сокращенный URL не найден", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
