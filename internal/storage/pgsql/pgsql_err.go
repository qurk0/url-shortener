package pgsql

import (
	"context"
	"errors"
	"net"
	"taskService/internal/lib/service/errs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func errMapping(err error) error {
	// Не нашли строки в БД
	if errors.Is(err, pgx.ErrNoRows) {
		return &errs.DbError{
			Code:    errs.CodeDbNotFound,
			Message: "zero rows found",
		}
	}

	// Отмена операции
	if errors.Is(err, context.Canceled) {
		return &errs.DbError{
			Code:    errs.CodeDbCancelled,
			Message: "operation canceled",
		}
	}

	// Конец времени выполнения операции
	if errors.Is(err, context.DeadlineExceeded) {
		return &errs.DbError{
			Code:    errs.CodeDbTimeout,
			Message: "timeout exceeded",
		}
	}

	// Проверяем ошибки сети
	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return &errs.DbError{
				Code:    errs.CodeDbTimeout,
				Message: "network timeout",
			}
		}
		return &errs.DbError{
			Code:    errs.CodeDbTemporary,
			Message: "temporary network error",
		}
	}

	// Проверяем ошибку уникальности
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return &errs.DbError{
				Code:    errs.CodeDbDuplicateAlias,
				Message: "alias already exist",
			}
		}
	}

	// Если никуда выше не попали - это внутренняя ошибка
	// Ошибка типа DbError логируется и клиенту не возвращается
	// Чтобы была хоть какая-то информация о внутренней ошибке, мы приплюсовываем содержимое ошибки
	return &errs.DbError{
		Code:    errs.CodeDbInternal,
		Message: "Internal storage error" + " " + err.Error(),
	}
}

func zeroRowsError() error {
	return &errs.DbError{
		Code:    errs.CodeDbNotFound,
		Message: "zero rows found",
	}
}
