package errs

type DbErrCode string

const (
	CodeDbNotFound       DbErrCode = "NOT_FOUND"
	CodeDbDuplicateAlias DbErrCode = "ALIAS_ALREADY_EXIST"
	CodeDbInternal       DbErrCode = "INTERNAL"
)

type DbError struct {
	Code    DbErrCode
	Message string
}

func (e *DbError) Error() string {
	return e.Message
}

type ServErrCode string

const (
	CodeServNotFound      ServErrCode = "NOT_FOUND"
	CodeServBadRequest    ServErrCode = "BAD_REQUEST"
	CodeServInternal      ServErrCode = "INTERNAL"
	CodeServAlreadyExists ServErrCode = "ALREADY_EXISTS"
	// CodeServUnauthorized ServErrCode = "UNAUTHORIZED" - задел на возможное будущее с миддлваром на аутентификацию
	// CodeServForbidden    ServErrCode = "FORBIDDEN" - задел на возможное будущее с добавлением прав доступа
)

type ServError struct {
	Code    ServErrCode
	Message string
}

func (e *ServError) Error() string {
	return e.Message
}
