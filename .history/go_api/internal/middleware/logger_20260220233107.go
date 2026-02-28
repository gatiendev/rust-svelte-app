package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
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
			Dur("duration", latency).
			Int("size", size).
			Str("user_agent", c.Request.UserAgent()).
			Msg("HTTP request")
	}
}
