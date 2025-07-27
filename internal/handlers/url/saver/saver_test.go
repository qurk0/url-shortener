package saver_test

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"taskService/internal/handlers/url/deleter"
	"taskService/internal/handlers/url/deleter/mocks"
	"taskService/internal/lib/service/errs"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveHandler_cases(t *testing.T) {
	cases := []struct {
		name           string
		alias          string
		mockReturnErr  error
		mockReturnID   int64
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success",
			alias:          "someAlias",
			mockReturnErr:  nil,
			mockReturnID:   12345,
			expectedStatus: fiber.StatusOK,
			expectedBody:   `{"status":"OK"}`,
		},
		{
			name:           "not_found",
			alias:          "someAlias",
			mockReturnErr:  &errs.DbError{Code: errs.CodeDbNotFound},
			expectedStatus: fiber.StatusNotFound,
			expectedBody: `
			{
				"status": "ERROR",
				"error": {
					"code": "NOT_FOUND",
					"message": "Your alias not found"
				}
			}
			`,
		},
		{
			name:           "internal",
			alias:          "someAlias",
			mockReturnErr:  &errs.DbError{Code: errs.CodeDbInternal},
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody: `
			{
				"status": "ERROR",
				"error": {
					"code": "INTERNAL",
					"message": "Somethings wrong in service. Try again later"
				}
			}
			`,
		},
		{
			name:           "temporary",
			alias:          "someAlias",
			mockReturnErr:  &errs.DbError{Code: errs.CodeDbTemporary},
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody: `
			{
				"status": "ERROR",
				"error": {
					"code": "TEMPORARY",
					"message": "Service temporary unavailable. Try again later"
				}
			}
			`,
		},
		{
			name:           "timeout",
			alias:          "someAlias",
			mockReturnErr:  &errs.DbError{Code: errs.CodeDbTimeout},
			expectedStatus: fiber.StatusRequestTimeout,
			expectedBody: `
			{
				"status": "ERROR",
				"error": {
					"code": "TIMEOUT",
					"message": "The server took too long to respond, try again"
				}
			}
			`,
		},
		{
			name:           "cancelled",
			alias:          "someAlias",
			mockReturnErr:  &errs.DbError{Code: errs.CodeDbCancelled},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: `
			{
				"status": "ERROR",
				"error": {
					"code": "CANCELLED",
					"message": "Operation cancelled"
				}
			}
			`,
		},
		{
			name:           "unexpected",
			alias:          "someAlias",
			mockReturnErr:  errors.New("Unexpected error"),
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody: `
			{
				"status": "ERROR",
				"error": {
					"code": "INTERNAL",
					"message": "Unexpected error"
				}
			}
			`,
		},
	}

	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockDeleter := mocks.NewURLDeleter(t)
			mockDeleter.On("DeleteURL", mock.Anything, tt.alias).Return(tt.mockReturnErr)

			app := fiber.New()
			app.Delete("/:alias", deleter.New(log, mockDeleter))

			req := httptest.NewRequest(http.MethodDelete, "/"+tt.alias, nil)

			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			assert.JSONEq(t, tt.expectedBody, string(body))

			mockDeleter.AssertCalled(t, "DeleteURL", mock.Anything, tt.alias)
		})
	}
}
