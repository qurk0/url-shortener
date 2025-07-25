package middleware

import "github.com/gofiber/fiber/v2"

func Validator[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body T
		if err := c.BodyParser(&body); err != nil {
			return dto.
		}

		return c.Next()
	}
}
