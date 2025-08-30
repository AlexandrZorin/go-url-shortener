package main

import (
	"log"

	"github.com/AlexandrZorin/go-url-shortener/internal/config"
	"github.com/AlexandrZorin/go-url-shortener/internal/handler"
	"github.com/AlexandrZorin/go-url-shortener/internal/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.CreateConfig()

	r := gin.Default()
	r.Use(logger.LoggerMiddleware())
	handler.SetupRoutes(r, cfg)
	r.Run(cfg.ServerAddress)
	return nil
}
