package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Logger is a middleware function that logs incoming requests and responses.
func Logger(c *fiber.Ctx) error {
	start := time.Now()

	// Process request
	err := c.Next()

	// Log request details
	log.Printf(
		"%s %s %s %v",
		c.Method(),
		c.OriginalURL(),
		c.Response().StatusCode(),
		time.Since(start),
	)

	return err
}
