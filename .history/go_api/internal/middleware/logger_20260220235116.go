package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for health checks
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		size := c.Writer.Size()

		log.Info().
			Str("method", method).
			Str("path", path).
			Str("remote_addr", clientIP).
			Int("status", statusCode).
			Str("duration", fmt.Sprintf("%dms", latency.Milliseconds())).
			Int("size", size).
			Str("user_agent", c.Request.UserAgent()).
			Msg("HTTP request")
	}
}
