package storage

import (
	"context"
	"errors"
	"log/slog"
	"taskService/internal/lib/service/errs"
	"taskService/internal/storage/pgsql"
	"taskService/internal/storage/redis"
	"time"
)

type Storage struct {
	mainDb *pgsql.Storage
	cache  *redis.Storage
	log    *slog.Logger
}

func New(mainDb *pgsql.Storage, cache *redis.Storage, log *slog.Logger) *Storage {
	return &Storage{
		mainDb: mainDb,
		cache:  cache,
		log:    log,
	}
}

func (s *Storage) SaveURL(ctx context.Context, url, alias string) (int64, error) {
	return s.mainDb.SaveURL(ctx, url, alias) // Просто вызываем метод главной БД
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	urlToRedirect, err := s.cache.GetURL(ctx, alias)
	if err != nil {
		var dbErr *errs.DbError
		if errors.As(err, &dbErr) {
			switch dbErr.Code {
			case errs.CodeDbCanceled, errs.CodeDbTimeout:
				return "", err
			case errs.CodeDbNotFound:
				urlToRedirect, err := s.mainDb.GetURL(ctx, alias)
				// Если ошибка прилетела из основной БД - она приведена в DbErr
				// Хендлеры рассчитаны на получение DbErr и сами приведут ошибку к ServError
				if err != nil {
					return "", err
				}

				// Асинхронно сохраняем в кеш наше соотношение, чтобы клиент не ждал выполнение этой операции
				go func(alias, url string) {
					internalCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					defer cancel()

					err := s.cache.SaveURL(internalCtx, urlToRedirect, alias)
					if err != nil {
						s.log.Error("failed to save url to cache", slog.String("error", err.Error()))
					}
				}(alias, urlToRedirect)
			default:
				s.log.Error("cache error",
					slog.String("alias", alias),
					slog.Any("error", err),
				)

				urlToRedirect, err = s.mainDb.GetURL(ctx, alias)
				if err != nil {
					return "", err
				}
			}
		} else {
			s.log.Error("unknown cache error", slog.Any("error", err))

			urlToRedirect, err = s.mainDb.GetURL(ctx, alias)
			if err != nil {
				return "", err
			}
		}
	}

	return urlToRedirect, nil
}

func (s *Storage) DeleteURL(ctx context.Context, alias string) error {
	err := s.mainDb.DeleteURL(ctx, alias)
	if err != nil {
		return err
	}

	err = s.cache.DeleteURL(ctx, alias)
	if err != nil {
		s.log.Error("failed to delete url from cache", slog.Any("error", err))
	}

	return nil
}
