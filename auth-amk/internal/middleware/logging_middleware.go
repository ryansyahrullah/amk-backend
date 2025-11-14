package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware memberikan log ringkas setiap request.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		status := c.Writer.Status()
		latency := time.Since(start)
		log.Printf("[GIN] %s %s -> %d (%s)", method, path, status, latency)
	}
}
