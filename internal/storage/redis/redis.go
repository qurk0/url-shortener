package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/qurk0/url-shortener/internal/config"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	client *redis.Client
	ttl    time.Duration
}

func New(cfg config.RedisConfig) (*Storage, error) {
	const op = "storage.redis.New"

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%s: redis ping failed: %w", op, err)
	}

	return &Storage{
		client: client,
		ttl:    cfg.TTL,
	}, nil
}

func (s *Storage) SaveURL(ctx context.Context, url, alias string) error {
	const op = "storage.pgsql.SaveURL"
	err := s.client.SetNX(ctx, alias, url, s.ttl).Err()
	if err != nil {
		err = errMapping(err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	val, err := s.client.GetEx(ctx, alias, s.ttl).Result()
	if err != nil {
		err = errMapping(err)
		return "", err
	}
	return val, nil
}

func (s *Storage) DeleteURL(ctx context.Context, alias string) error {
	_, err := s.client.Del(ctx, alias).Result()

	if err != nil {
		err = errMapping(err)
		return err
	}

	return nil
}
