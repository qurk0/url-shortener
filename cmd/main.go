package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"taskService/internal/config"
	"taskService/internal/lib/log/sl"
	"taskService/internal/storage/pgsql"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	log := setupLogger(cfg.Env)
	log.Info("Starting task service...", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled")

	storage, err := pgsql.New(context.Background(), cfg.DbCfg)
	if err != nil {
		log.Error("failed to create storage", sl.Err(err))
		os.Exit(1)
	}

	// TODO: Начало обслуживания адреса

	// TODO: Шатдаун
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
