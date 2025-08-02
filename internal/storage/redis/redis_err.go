package redis

import (
	"context"
	"errors"
	"net"

	"github.com/qurk0/url-shortener/internal/lib/service/errs"

	"github.com/redis/go-redis/v9"
)

func errMapping(err error) error {
	// В редисе нет ключа
	if errors.Is(err, redis.Nil) {
		return &errs.DbError{
			Code:    errs.CodeDbNotFound,
			Message: "Key not exists in Redis",
		}
	}

	// Отмена контекста
	if errors.Is(err, context.Canceled) {
		return &errs.DbError{
			Code:    errs.CodeDbCancelled,
			Message: "Operation canceled",
		}
	}

	// Окончание времени выполнения операции
	if errors.Is(err, context.DeadlineExceeded) {
		return &errs.DbError{
			Code:    errs.CodeDbTimeout,
			Message: "Timeout reached",
		}
	}

	// Сетевые ошибки
	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return &errs.DbError{
				Code:    errs.CodeDbTimeout,
				Message: "Timeout reached",
			}
		}

		return &errs.DbError{
			Code:    errs.CodeDbTemporary,
			Message: "Temporary network error",
		}
	}

	// Если никуда выше не попали - это внутренняя ошибка
	// Ошибка типа DbError логируется и клиенту не возвращается
	// Чтобы была хоть какая-то информация о внутренней ошибке, мы приплюсовываем содержимое ошибки
	return &errs.DbError{
		Code:    errs.CodeDbInternal,
		Message: "Internal error" + " " + err.Error(),
	}
}
