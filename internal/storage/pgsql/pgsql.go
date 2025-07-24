package pgsql

import (
	"context"
	"fmt"
	"taskService/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

const (
	SaveURLQuery   = "INSERT INTO url (url, alias) VALUES ($1, $2) RETURNING id"
	GetURLQuery    = "SELECT url FROM url WHERE alias = $1"
	DeleteURLQuery = "DELETE FROM url WHERE alias = $1"
)

func New(ctx context.Context, cfg config.DBConfig) (*Storage, error) {
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
		return nil, fmt.Errorf("failed to parse PostgreSQL storage configs: %v", err)
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create PostgreSQL storage pool: %v", err)
	}

	return &Storage{pool: pool}, err
}

func (s *Storage) SaveURL(ctx context.Context, urlToSave, alias string) (int64, error) {
	const op = "storage.pgsql.SaveURL"

	var id int64
	err := s.pool.QueryRow(ctx, SaveURLQuery, urlToSave, alias).Scan(&id)
	if err != nil {
		err = errMapping(err)
		return -1, err
	}

	return id, nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "storage.pgsql.GetURL"

	var urlResp string
	err := s.pool.QueryRow(ctx, GetURLQuery, alias).Scan(&urlResp)
	if err != nil {
		err = errMapping(err)
		return "", err
	}

	return urlResp, nil
}

func (s *Storage) DeleteURL(ctx context.Context, alias string) error {
	const op = "storage.pgsql.DeleteURL"

	tags, err := s.pool.Exec(ctx, DeleteURLQuery, alias)
	if err != nil {
		err = errMapping(err)
		return err
	}

	if tags.RowsAffected() == 0 {
		err = errMapping(pgx.ErrNoRows)
		return err
	}

	return nil
}
