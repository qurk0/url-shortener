package saver

import (
	"context"
	"crypto/rand"
	"errors"
	"log/slog"
	"math/big"

	resp "github.com/qurk0/url-shortener/internal/lib/service/api/response"
	"github.com/qurk0/url-shortener/internal/lib/service/errs"

	"github.com/gofiber/fiber/v2"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLSaver --output=./mocks --outpkg=mocks
type URLSaver interface {
	SaveURL(ctx context.Context, url, alias string, userID int) (int64, error)
}

const (
	randomAliasLength = 10

	lowerCaseSymbols = "abcdefghijklmnopqrstuvwxyz"
	upperCaseSymbols = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers          = "0123456789"
	specialSymbols   = "?!@*"

	alphabet = lowerCaseSymbols + upperCaseSymbols + numbers
)

func New(log *slog.Logger, urlSaver URLSaver) fiber.Handler {
	return func(c *fiber.Ctx) error {
		const op = "handlers.url.save"

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

		log.Debug("Starting url-alias creation")

		userIDRaw := c.Locals("user_id")
		userID, ok := userIDRaw.(int)
		if !ok {
			log.Error("invalid user_id in request context", slog.Any("value", userIDRaw))
			return resp.ReturnError(c, errs.ServError{
				Code:    errs.CodeServInternal,
				Message: "Invalid User ID",
			})
		}

		reqRaw := c.Locals("validated-body")
		req, ok := reqRaw.(Request)
		if !ok {
			log.Error("invalid validated-body type", slog.Any("value", reqRaw))
			return resp.ReturnError(c, errs.ServError{
				Code:    errs.CodeServInternal,
				Message: "Invalid request body format",
			})
		}

		if req.Alias == "" {
			// Заменяем пустой алиас на сгенерированный
			newAlias, err := generateAlias()
			if err != nil {
				log.Error("alias generation error",
					slog.Any("error", err))
				return resp.ReturnError(c, errs.ServError{
					Code:    errs.CodeServInternal,
					Message: "Failed to generate new alias",
				})
			}

			req.Alias = newAlias
		}

		log.Debug("Starting url-alias creation in db")

		id, err := urlSaver.SaveURL(c.Context(), req.URL, req.Alias, userID)
		if err != nil {
			log.Error("error from storage", slog.Any("error", err))
			var dbErr *errs.DbError
			if errors.As(err, &dbErr) {
				servErr := errs.ServError{
					Code: errs.MappingDbToServErrs(dbErr.Code),
				}

				switch servErr.Code {
				case errs.CodeServConflict:
					servErr.Message = "Your alias already exists, choose another one"
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
					slog.Any("alias", req.Alias),
				)
				return resp.ReturnError(c, servErr)
			}

			log.Error("unmapped error",
				slog.Any("alias", req.Alias),
				slog.Any("error", err),
			)

			return resp.ReturnError(c, errs.ServError{
				Code:    errs.CodeServInternal,
				Message: "Unexpected error",
			})
		}

		return resp.ReturnCreated(c, fiber.Map{"id": id})
	}
}

func generateAlias() (string, error) {
	alias := make([]byte, randomAliasLength)
	alphabetLen := big.NewInt(int64(len(alphabet)))
	for i := range alias {
		randNum, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}

		alias[i] = alphabet[randNum.Int64()]
	}

	return string(alias), nil
}
