package deleter

import (
	"context"
	"errors"
	"log/slog"
	authgrpc "taskService/internal/client/auth/grpc"
	resp "taskService/internal/lib/service/api/response"
	"taskService/internal/lib/service/errs"

	"github.com/gofiber/fiber/v2"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLDeleter --output=./mocks --outpkg=mocks
type URLDeleter interface {
	DeleteURL(ctx context.Context, alias string) error
	IsOwner(ctx context.Context, alias string, userID int) (bool, error)
}

func New(log *slog.Logger, urlDeleter URLDeleter, authClient *authgrpc.Client) fiber.Handler {
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

		userIDRaw := c.Locals("user_id")
		userID, ok := userIDRaw.(int)
		if !ok {
			log.Warn("unexpected missing or invalid user_id", slog.Any("value", userIDRaw))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal auth error",
			})
		}

		isOwner, err := urlDeleter.IsOwner(c.Context(), alias, userID)
		if err != nil {
			log.Error("error from storage", slog.Any("error", err))
			return ErrReturn(c, err, alias)
		}

		if !isOwner {
			isAdmin, err := authClient.IsAdmin(c.Context(), int64(userID))
			if err != nil {
				log.Error("error from gRPC client", slog.Any("error", err))
				return resp.ReturnError(c, errs.ServError{Code: errs.CodeServInternal, Message: "Internal error, try again later"})
			}

			if !isAdmin {
				log.Info("Invalid attempt to delete URL Alias", slog.Int("userID", userID))
				return resp.ReturnForbidden(c)
			}
		}

		// До сюда мы дойдет ТОЛЬКО при одном из условий
		// - мы - владелец алиаса
		// - мы - администратор

		err = urlDeleter.DeleteURL(c.Context(), alias)
		if err != nil {
			log.Error("error from storage", slog.Any("error", err))
			return ErrReturn(c, err, alias)
		}

		return resp.ReturnOk(c)
	}
}

func ErrReturn(c *fiber.Ctx, err error, alias string) error {
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
		case errs.CodeServTemporary:
			servErr.Message = "Service temporary unavailable. Try again later"
		case errs.CodeServTimeout:
			servErr.Message = "The server took too long to respond, try again"
		case errs.CodeServCancelled:
			servErr.Message = "Operation cancelled"
		}
		return resp.ReturnError(c, servErr)
	}
	return resp.ReturnError(c, errs.ServError{
		Code:    errs.CodeServInternal,
		Message: "Unexpected error",
	})
}
