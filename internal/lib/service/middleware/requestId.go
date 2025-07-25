package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := c.Get(RequestIDHeader)
		if reqID == "" {
			reqID = uuid.NewString()
		}

		c.Locals(RequestIDHeader, reqID)
		c.Set(RequestIDHeader, reqID)

		return c.Next()
	}
}
