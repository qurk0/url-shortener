package deleter

import (
	"context"
	"errors"
	"log/slog"
	resp "taskService/internal/lib/service/api/response"
	"taskService/internal/lib/service/errs"

	"github.com/gofiber/fiber/v2"
)

type URLDeleter interface {
	DeleteURL(ctx context.Context, alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		const op = "handlers.url.delete"

		reqIDRaw := c.Locals("X-Request-ID")
		reqID, ok := reqIDRaw.(string)
		if !ok {
			reqID = "Unknown"
			log.Warn("missing or invalid request id", slog.Any("value", reqIDRaw))
		}

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", reqID),
		)

		alias := c.Params("alias")
		if alias == "" {
			return resp.ReturnError(c, errs.ServError{
				Code:    errs.CodeServBadRequest,
				Message: "Empty alias in your URL",
			})
		}

		err := urlDeleter.DeleteURL(c.Context(), alias)
		if err != nil {
			log.Error("error from storage", slog.Any("error", err))
			var dbErr *errs.DbError
			if errors.As(err, &dbErr) {
				servErr := errs.ServError{
					Code: errs.MappingDbToServErrs((dbErr.Code)),
				}

				switch servErr.Code {
				case errs.CodeServNotFound:
					servErr.Message = "Your alias not found"
				case errs.CodeServInternal:
					servErr.Message = "Somethings wrong in service. Try again later"
				default:
					servErr.Message = "Unknown error. Write to our support for help"
				}
				log.Error("mapped service error",
					slog.Any("error", servErr),
					slog.Any("alias", alias),
				)
				return resp.ReturnError(c, servErr)
			}

			log.Error("unmapped error",
				slog.Any("error", err),
				slog.Any("alias", alias),
			)
			return resp.ReturnError(c, errs.ServError{
				Code:    errs.CodeServInternal,
				Message: "Unexpected error",
			})
		}

		return resp.ReturnOk(c)
	}
}
