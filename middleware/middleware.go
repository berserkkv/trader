package middleware

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

func SlogLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		latency := time.Since(start)

		status := c.Writer.Status()

		slog.Info("HTTP request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", status,
			"latency", latency,
			"client_ip", c.ClientIP(),
		)
	}
}
