package pgsql

import (
	"errors"
	"taskService/internal/lib/service/errs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func errMapping(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return &errs.DbError{
			Code:    errs.CodeNotFound,
			Message: "zero rows found",
		}
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return &errs.DbError{
				Code:    errs.ErrDuplicateAlias,
				Message: "alias already exist",
			}
		}
	}

	return &errs.DbError{
		Code:    errs.ErrInternal,
		Message: "Internal storage error",
	}
}
