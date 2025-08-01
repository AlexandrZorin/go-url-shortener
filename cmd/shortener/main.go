package main

import (
	"net/http"
	"strings"
)

type URL struct {
	urls map[string]string
}

var urlStorage = URL{
	urls: make(map[string]string),
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, postHandler)
	mux.HandleFunc(`/{id}`, getHandler)
	err := http.ListenAndServe(`:8080`, nil)
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
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Этот URL принимает только GET запросы", http.StatusBadRequest)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
}
