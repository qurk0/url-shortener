package middleware

import (
	resp "taskService/internal/lib/service/api/response"
	"taskService/internal/lib/service/errs"
	"taskService/internal/lib/service/validator"

	"github.com/gofiber/fiber/v2"
)

func Validator[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body T
		if err := c.BodyParser(&body); err != nil {
			return resp.ReturnError(c, errs.ServError{
				Code:    errs.CodeServBadRequest,
				Message: "Failed to parse request body",
			})
		}

		if err := validator.Validate(c.Context(), body); err != nil {
			return resp.ReturnError(c, errs.ServError{
				Code:    errs.CodeServBadRequest,
				Message: "Failed to validate request body",
			})
		}

		c.Locals("validated-body", body)

		return c.Next()
	}
}
