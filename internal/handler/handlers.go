package handler

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

type URLStorage struct {
	urls map[string]string
}

var storage = URLStorage{
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

func SetupRoutes(r *gin.Engine) {
	r.POST("/", handlePostRequest)
	r.GET("/:id", handleGetRequest)
}

func handlePostRequest(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.String(http.StatusBadRequest, "Этот URL принимает только POST запросы")
		return
	}
	if c.ContentType() != "text/plain" {
		c.String(http.StatusBadRequest, "Неверный Content-Type. Ожидается text/plain")
		return
	}
	originalURL, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusBadRequest, "Ошибка чтения тела запроса")
		return
	}
	if len(originalURL) == 0 {
		c.String(http.StatusBadRequest, "URL не может быть пустым")
		return
	}
	var shortKey string
	for k, v := range storage.urls {
		if v == string(originalURL) {
			shortKey = k
			break
		}
	}
	if shortKey == "" {
		shortKey = generateShortURL()
		storage.urls[shortKey] = string(originalURL)
	}
	shortURL := fmt.Sprintf("http://localhost:8080/%s", shortKey)
	c.Data(http.StatusCreated, "text/plain", []byte(shortURL))
}

func handleGetRequest(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.String(http.StatusBadRequest, "Этот URL принимает только GET запросы")
		return
	}
	shortKey := c.Param("id")
	if shortKey == "" {
		c.String(http.StatusBadRequest, "Неверный запрос")
		return
	}
	originalURL, exists := storage.urls[shortKey]
	if !exists {
		c.String(http.StatusBadRequest, "Сокращенный URL не найден")
		return
	}
	c.Header("Location", originalURL)
	c.String(http.StatusTemporaryRedirect, "")
}
