package logger

import (
	"time"

	"github.com/gin-gonic/gin"
)

// GinLogger Gin中间件日志
func GinLogger(log Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Info("HTTP Request",
			"status", statusCode,
			"latency", latency,
			"client_ip", clientIP,
			"method", method,
			"path", path,
		)
	}
}
