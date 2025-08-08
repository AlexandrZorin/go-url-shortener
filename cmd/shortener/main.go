package main

import (
	"github.com/AlexandrZorin/go-url-shortener/internal/config"
	"github.com/AlexandrZorin/go-url-shortener/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.CreateConfig()

	r := gin.Default()
	handler.SetupRoutes(r, cfg)
	r.Run(":8080")
}
