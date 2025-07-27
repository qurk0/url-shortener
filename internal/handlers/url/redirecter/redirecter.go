package redirecter

import (
	"context"
	"errors"
	"log/slog"
	resp "taskService/internal/lib/service/api/response"
	"taskService/internal/lib/service/errs"

	"github.com/gofiber/fiber/v2"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLGetter --output=./mocks --outpkg=mocks
type URLGetter interface {
	GetURL(ctx context.Context, alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		const op = "handlers.url.redirect"

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

		url, err := urlGetter.GetURL(c.Context(), alias)
		if err != nil {
			log.Error("error from storage", slog.Any("error", err))
			var dbErr *errs.DbError
			if errors.As(err, &dbErr) {
				servErr := errs.ServError{
					Code: errs.MappingDbToServErrs(dbErr.Code),
				}

				switch servErr.Code {
				case errs.CodeServNotFound:
					servErr.Message = "Your alias not found"
				case errs.CodeServInternal:
					servErr.Message = "Somethings wrong in service. Try again later"
				case errs.CodeServTemporary:
					servErr.Message = "Service temporary unavailable. Try again later"
				case errs.CodeServTimeout:
					servErr.Message = "The server took too long to respond, try again"
				case errs.CodeServCancelled:
					servErr.Message = "Operation cancelled"
				}
				log.Error("mapped service error",
					slog.Any("error", servErr),
					slog.Any("alias", alias),
				)
				return resp.ReturnError(c, servErr)
			}

			log.Error("unmapped error",
				slog.Any("alias", alias),
				slog.Any("error", err),
			)

			return resp.ReturnError(c, errs.ServError{
				Code:    errs.CodeServInternal,
				Message: "Unexpected error",
			})
		}

		return resp.ReturnRedirecting(c, url)
	}
}
