package middleware

import (
	"time"

	"fitgenie/pkg/logger"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		timestamp := time.Now()
		latency := timestamp.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Info("request",
			"timestamp", timestamp.Format(time.RFC3339),
			"latency", latency,
			"client_ip", clientIP,
			"method", method,
			"path", path,
			"status", statusCode,
			"errors", c.Errors.ByType(gin.ErrorTypePrivate).String(),
		)
	}
}
