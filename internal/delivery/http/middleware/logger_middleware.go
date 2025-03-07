//internal/delivery/http/middleware/logger_middleware.go

package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LoggerMiddleware struct{}

func NewLoggerMiddleware() *LoggerMiddleware {
	return &LoggerMiddleware{}
}

func (m *LoggerMiddleware) LogRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		
		// Proses request
		err := c.Next()
		
		// Log request
		duration := time.Since(start)
		log.Printf(
			"[%s] %s %s %d %s",
			c.Method(),
			c.Path(),
			c.IP(),
			c.Response().StatusCode(),
			duration,
		)
		
		return err
	}
}