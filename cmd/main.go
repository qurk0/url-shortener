package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"taskService/internal/config"
	"taskService/internal/handlers/url/deleter"
	"taskService/internal/handlers/url/redirecter"
	"taskService/internal/handlers/url/saver"
	"taskService/internal/lib/log/sl"
	"taskService/internal/lib/service/middleware"
	"taskService/internal/storage"
	"taskService/internal/storage/pgsql"
	"taskService/internal/storage/redis"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Получение конфигов
	cfg := config.MustLoad()

	// Инициализация логгера
	log := setupLogger(cfg.Env)
	log.Info("Starting task service...", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled")

	// Инициализация стореджа

	mainStorage, err := pgsql.New(context.Background(), cfg.MainDBCfg)
	if err != nil {
		log.Error("failed to create main storage", sl.Err(err))
		os.Exit(1)
	}

	cacheStorage, err := redis.New(cfg.CacheCfg)
	if err != nil {
		log.Error("failed to create cache storage", sl.Err(err))
		os.Exit(1)
	}

	storage := storage.New(mainStorage, cacheStorage, log)

	// TODO: Начало обслуживания адреса
	app := fiber.New()

	app.Post("/url/new",
		middleware.RequestID(),
		middleware.Logger(log),
		middleware.Validator[saver.Request](),
		saver.New(log, storage),
	)

	app.Get("/:alias",
		middleware.RequestID(),
		middleware.Logger(log),
		redirecter.New(log, storage),
	)

	app.Delete("/:alias",
		middleware.RequestID(),
		middleware.Logger(log),
		deleter.New(log, storage),
	)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info("Starting listen URL Shortener",
			slog.String("addr", cfg.ApiCfg.Addr),
		)
		if err := app.Listen(cfg.ApiCfg.Addr); err != nil {
			log.Error("URL Shortener didn't start",
				slog.Any("error", err),
			)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down URL Shortener")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Error("Shutdown failed",
			slog.Any("error", err),
		)
	} else {
		log.Info("URL Shortener stopped gracefully")
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return log
}
