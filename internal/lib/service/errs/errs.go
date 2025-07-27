package errs

// Коды ошибки на уровне репозитория
// Они рассчитаны, что они будут возвращены из репозитория в сервис и будут использованы для логирования
// И далее будут преобразованы в ServError с помощью маппинга
type DbErrCode string

const (
	CodeDbNotFound       DbErrCode = "NOT_FOUND"
	CodeDbDuplicateAlias DbErrCode = "ALIAS_ALREADY_EXIST"
	CodeDbInternal       DbErrCode = "INTERNAL"
	CodeDbCancelled      DbErrCode = "CANCELLED"
	CodeDbTimeout        DbErrCode = "TIMEOUT"
	CodeDbTemporary      DbErrCode = "TEMPORARY"
)

type DbError struct {
	Code    DbErrCode
	Message string
}

func (e *DbError) Error() string {
	return e.Message
}

// Коды ошибки на уровне сервиса
// Они рассчитаны, что они полетят к клиенту в теле ответа
type ServErrCode string

const (
	// Не мапим
	CodeServBadRequest ServErrCode = "BAD_REQUEST"

	// Мапим
	CodeServNotFound  ServErrCode = "NOT_FOUND"
	CodeServInternal  ServErrCode = "INTERNAL"
	CodeServConflict  ServErrCode = "CONFLICT"
	CodeServCancelled ServErrCode = "CANCELLED"
	CodeServTimeout   ServErrCode = "TIMEOUT"
	CodeServTemporary ServErrCode = "TEMPORARY"
	// CodeServUnauthorized ServErrCode = "UNAUTHORIZED" - задел на возможное будущее с миддлваром на аутентификацию
	// CodeServForbidden    ServErrCode = "FORBIDDEN" - задел на возможное будущее с добавлением прав доступа
)

type ServError struct {
	Code    ServErrCode `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
}

func (e *ServError) Error() string {
	return e.Message
}

// Маппинг ошибок из БД к ошибкам сервиса

var mapping map[DbErrCode]ServErrCode = map[DbErrCode]ServErrCode{
	CodeDbNotFound:       CodeServNotFound,
	CodeDbDuplicateAlias: CodeServConflict,
	CodeDbInternal:       CodeServInternal,
	CodeDbCancelled:      CodeServCancelled,
	CodeDbTimeout:        CodeServTimeout,
	CodeDbTemporary:      CodeServTemporary,
}

func MappingDbToServErrs(code DbErrCode) ServErrCode {
	return mapping[code]
}
