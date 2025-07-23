package pgsql

import (
	"taskService/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(cfg config.DBConfig) (*Storage, error) {

}
