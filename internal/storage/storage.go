package storage

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrExists   = errors.New("already exists")
	ErrInternal = errors.New("internal error")
)
