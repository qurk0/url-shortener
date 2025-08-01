package redirecter_test

// import (
// 	"errors"
// 	"io"
// 	"log/slog"
// 	"net/http"
// 	"net/http/httptest"
// 	"taskService/internal/handlers/url/redirecter"
// 	"taskService/internal/handlers/url/redirecter/mocks"
// 	"taskService/internal/lib/service/errs"
// 	"testing"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestRedirectHandler_cases(t *testing.T) {
// 	cases := []struct {
// 		name           string
// 		alias          string
// 		mockReturnErr  error
// 		mockReturnUrl  string
// 		expectedStatus int
// 		expectedBody   string
// 	}{
// 		{
// 			name:           "success",
// 			alias:          "someAlias",
// 			mockReturnErr:  nil,
// 			mockReturnUrl:  "http://example.com",
// 			expectedStatus: fiber.StatusFound,
// 			expectedBody:   "",
// 		},
// 		{
// 			name:           "not_found",
// 			alias:          "someAlias",
// 			mockReturnErr:  &errs.DbError{Code: errs.CodeDbNotFound},
// 			mockReturnUrl:  "",
// 			expectedStatus: fiber.StatusNotFound,
// 			expectedBody: `
// 			{
// 				"status": "ERROR",
// 				"error": {
// 					"code": "NOT_FOUND",
// 					"message": "Your alias not found"
// 				}
// 			}
// 			`,
// 		},
// 		{
// 			name:           "internal",
// 			alias:          "someAlias",
// 			mockReturnErr:  &errs.DbError{Code: errs.CodeDbInternal},
// 			expectedStatus: fiber.StatusInternalServerError,
// 			expectedBody: `
// 			{
// 				"status": "ERROR",
// 				"error": {
// 					"code": "INTERNAL",
// 					"message": "Somethings wrong in service. Try again later"
// 				}
// 			}
// 			`,
// 		},
// 		{
// 			name:           "temporary",
// 			alias:          "someAlias",
// 			mockReturnErr:  &errs.DbError{Code: errs.CodeDbTemporary},
// 			mockReturnUrl:  "",
// 			expectedStatus: fiber.StatusInternalServerError,
// 			expectedBody: `
// 			{
// 				"status": "ERROR",
// 				"error": {
// 					"code": "TEMPORARY",
// 					"message": "Service temporary unavailable. Try again later"
// 				}
// 			}
// 			`,
// 		},
// 		{
// 			name:           "timeout",
// 			alias:          "someAlias",
// 			mockReturnErr:  &errs.DbError{Code: errs.CodeDbTimeout},
// 			mockReturnUrl:  "",
// 			expectedStatus: fiber.StatusRequestTimeout,
// 			expectedBody: `
// 			{
// 				"status": "ERROR",
// 				"error": {
// 					"code": "TIMEOUT",
// 					"message": "The server took too long to respond, try again"
// 				}
// 			}
// 			`,
// 		},
// 		{
// 			name:           "cancelled",
// 			alias:          "someAlias",
// 			mockReturnErr:  &errs.DbError{Code: errs.CodeDbCancelled},
// 			mockReturnUrl:  "",
// 			expectedStatus: fiber.StatusBadRequest,
// 			expectedBody: `
// 			{
// 				"status": "ERROR",
// 				"error": {
// 					"code": "CANCELLED",
// 					"message": "Operation cancelled"
// 				}
// 			}
// 			`,
// 		},
// 		{
// 			name:           "unexpected",
// 			alias:          "someAlias",
// 			mockReturnErr:  errors.New("Unexpected error"),
// 			mockReturnUrl:  "",
// 			expectedStatus: fiber.StatusInternalServerError,
// 			expectedBody: `
// 			{
// 				"status": "ERROR",
// 				"error": {
// 					"code": "INTERNAL",
// 					"message": "Unexpected error"
// 				}
// 			}
// 			`,
// 		},
// 	}

// 	log := slog.New(slog.NewTextHandler(io.Discard, nil))

// 	for _, tt := range cases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			mockGetter := mocks.NewURLGetter(t)
// 			mockGetter.On("GetURL", mock.Anything, tt.alias).Return(tt.mockReturnUrl, tt.mockReturnErr)

// 			app := fiber.New()
// 			app.Get("/:alias", redirecter.New(log, mockGetter))

// 			req := httptest.NewRequest(http.MethodGet, "/"+tt.alias, nil)

// 			resp, err := app.Test(req)

// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

// 			if tt.expectedBody != "" {
// 				body, _ := io.ReadAll(resp.Body)
// 				assert.JSONEq(t, tt.expectedBody, string(body))
// 			}

// 			mockGetter.AssertCalled(t, "GetURL", mock.Anything, tt.alias)
// 		})
// 	}
// }
