package response

import (
	"taskService/internal/lib/service/errs"

	"github.com/gofiber/fiber/v2"
)

type status string

const (
	StatusOk          status = "OK"
	StatusCreated     status = "CREATED"
	StatusRedirecting status = "REDIRECTING"
	StatusError       status = "ERROR"
)

type Response struct {
	Status status          `json:"status"`
	Error  *errs.ServError `json:"error,omitempty"`
	Data   any             `json:"data,omitempty"`
}

func ReturnError(ctx *fiber.Ctx, err errs.ServError) error {
	return ctx.Status(parseErrorToStatus(err.Code)).JSON(Response{
		Status: StatusError,
		Error:  &err,
	})
}

func ReturnCreated(ctx *fiber.Ctx, data any) error {
	return ctx.Status(fiber.StatusCreated).JSON(Response{
		Status: StatusCreated,
		Data:   data,
	})
}

func ReturnOk(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(Response{
		Status: StatusOk,
	})
}

func ReturnRedirecting(ctx *fiber.Ctx, redirectUrl string) error {
	// Выставляем заголовки, не дающие кэшировать ответ сервиса
	// ctx.Set("Cache-Control", "no-store, no-cache, must-revalidate")
	// ctx.Set("Pragma", "no-cache")
	// ctx.Set("Expires", "0")

	return ctx.Redirect(redirectUrl, fiber.StatusFound)
}

func parseErrorToStatus(code errs.ServErrCode) int {
	switch code {
	case errs.CodeServConflict:
		return fiber.StatusConflict
	case errs.CodeServBadRequest:
		return fiber.StatusBadRequest
	case errs.CodeServInternal:
		return fiber.StatusInternalServerError
	case errs.CodeServNotFound:
		return fiber.StatusNotFound
	// case CodeServUnauthorized:
	// 	return fiber.StatusUnauthorized
	// case CodeServForbidden:
	// 	return fiber.StatusForbidden
	default:
		return fiber.StatusInternalServerError
	}
}
