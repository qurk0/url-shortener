package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger(log *slog.Logger) func(c *fiber.Ctx) error {
	log = log.With(
		slog.String("component", "middleware/logger"),
	)

	log.Info("logger middleware enabled")
	return func(c *fiber.Ctx) error {
		start := time.Now()

		entry := log.With(
			slog.String("method", c.Method()),
			slog.String("path", c.OriginalURL()),
			slog.String("remote_addr", c.IP()),
			slog.String("user_agent", c.Get("User-Agent")),
			slog.String("request_id", c.Get("X-Request-ID")),
		)

		defer func() {
			entry.Info("request completed",
				slog.Int("status", c.Response().StatusCode()),
				slog.Int("bytes", len(c.Response().Body())),
				slog.String("duration", time.Since(start).String()),
			)
		}()

		return c.Next()
	}
}
