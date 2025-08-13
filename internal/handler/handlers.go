package handler

import (
	"fmt"
	"net/http"

	"github.com/AlexandrZorin/go-url-shortener/internal/config"
	"github.com/AlexandrZorin/go-url-shortener/internal/service"
	"github.com/gin-gonic/gin"
)

var URLService = service.NewURLService()

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	r.POST("/", func(c *gin.Context) {
		handlePostRequest(c, cfg)
	})
	r.GET("/:id", handleGetRequest)
}

func handlePostRequest(c *gin.Context, cfg *config.Config) {
	if c.Request.Method != http.MethodPost {
		c.String(http.StatusMethodNotAllowed, "Этот URL принимает только POST запросы")
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

	shortKey := URLService.CreateShortURL(string(originalURL))
	shortURL := fmt.Sprintf("%s/%s", cfg.URL, shortKey)
	c.Data(http.StatusCreated, "text/plain", []byte(shortURL))
}

func handleGetRequest(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.String(http.StatusMethodNotAllowed, "Этот URL принимает только GET запросы")
		return
	}
	shortKey := c.Param("id")
	if shortKey == "" {
		c.String(http.StatusBadRequest, "Неверный запрос")
		return
	}

	originalURL, exists := URLService.GetOriginalURL(shortKey)
	if !exists {
		c.String(http.StatusNotFound, "Сокращенный URL не найден")
		return
	}
	c.Header("Location", originalURL)
	c.String(http.StatusTemporaryRedirect, "")
}
