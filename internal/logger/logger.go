package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		logger.Info("Request/Response Info",
			zap.String("uri", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Duration("execution_time", duration),
			zap.Int("status_code", c.Writer.Status()),
			zap.Int("response_size", c.Writer.Size()),
		)
	}
}
