package errs

type Code string

const (
	CodeNotFound      = "NOT_FOUND"
	ErrDuplicateAlias = "ALIAS_ALREADY_EXIST"
	ErrInternal       = "INTERNAL"
)

type DbError struct {
	Code    Code
	Message string
}

func (e *DbError) Error() string {
	return e.Message
}
