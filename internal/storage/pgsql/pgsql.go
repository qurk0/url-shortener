package pgsql

import (
	"context"
	"fmt"

	"github.com/qurk0/url-shortener/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

const (
	SaveURLQuery          = "INSERT INTO urls (url, alias, user_id) VALUES ($1, $2, $3) RETURNING id"
	GetURLQuery           = "SELECT url FROM urls WHERE alias = $1"
	ChechUserIsOwnerQuery = "SELECT EXISTS (SELECT 1 FROM urls WHERE alias = $1 AND user_id = $2);"
	DeleteURLQuery        = "DELETE FROM urls WHERE alias = $1"
)

func New(ctx context.Context, cfg config.PGSQLConfig) (*Storage, error) {
	const op = "storage.pgsql.New"

	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
		cfg.SslMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime,
		cfg.PoolMaxConnIdleTime,
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse PostgreSQL storage configs: %v", op, err)
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create PostgreSQL storage pool: %v", op, err)
	}

	return &Storage{pool: pool}, err
}

func (s *Storage) SaveURL(ctx context.Context, urlToSave, alias string, userID int) (int64, error) {
	const op = "storage.pgsql.SaveURL"

	var id int64
	err := s.pool.QueryRow(ctx, SaveURLQuery, urlToSave, alias, userID).Scan(&id)
	if err != nil {
		err = errMapping(err)
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "storage.pgsql.GetURL"

	var urlResp string
	err := s.pool.QueryRow(ctx, GetURLQuery, alias).Scan(&urlResp)
	if err != nil {
		err = errMapping(err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return urlResp, nil
}

func (s *Storage) DeleteURL(ctx context.Context, alias string) error {
	const op = "storage.pgsql.DeleteURL"

	tags, err := s.pool.Exec(ctx, DeleteURLQuery, alias)
	if err != nil {
		err = errMapping(err)
		return fmt.Errorf("%s: %w", op, err)
	}

	if tags.RowsAffected() == 0 {
		err = zeroRowsError()
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) IsOwner(ctx context.Context, alias string, userID int) (bool, error) {
	const op = "storage.pgsql.IsOwner"

	var isUserOwner bool
	err := s.pool.QueryRow(ctx, ChechUserIsOwnerQuery, alias, userID).Scan(&isUserOwner)
	if err != nil {
		err = errMapping(err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isUserOwner, nil
}
