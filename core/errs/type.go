package errs

const (
	TypeIllegalArgument    Type = "ILLEGAL_ARGUMENT"    // Входной параметр или запрос не может быть принят.
	TypeAuthFailure        Type = "AUTH_FAILURE"        // Ошибка аутентификации.
	TypeForbidden          Type = "FORBIDDEN"           // Выполнение операции запрещено.
	TypeNotFound           Type = "NOT_FOUND"           // Объект не найден.
	TypeConflict           Type = "CONFLICT"            // Операция не согласуется с текущим состоянием.
	TypeCancelled          Type = "CANCELLED"           // Операция была отменена.
	TypeHasReferences      Type = "HAS_REFERENCES"      // Имеются ссылки на объект.
	TypeNotImplemented     Type = "NOT_IMPLEMENTED"     // Отсутствует реализация операции.
	TypeSystemFailure      Type = "SYSTEM_FAILURE"      // Не классифицируемая системная ошибка.
	TypeServiceUnavailable Type = "SERVICE_UNAVAILABLE" // Сервис недоступен.
	TypeTooManyRequests    Type = "TOO_MANY_REQUESTS"   // Превышено количество запросов.
	TypeUnauthenticated    Type = "UNAUTHENTICATED"     // Не авторизован
	TypePermissionDenied   Type = "PERMISSION_DENIED"   // Недостаточно прав для выполнения операции
	TypeResourceExhausted  Type = "RESOURCE_EXHAUSTED"  // Ресурс был исчерпан
	TypeFailedPrecondition Type = "FAILED_PRECONDITION" // Операциeя не может быть выполнена
	TypeAborted            Type = "ABORTED"             // Операция была прервана
	TypeOutOfRange         Type = "OUT_OF_RANGE"        // Операция не может быть выполнена
	TypeUnimplemented      Type = "UNIMPLEMENTED"       // Операция не реализована
	TypeInternal           Type = "INTERNAL"            // Внутренняя ошибка
	TypeUnavailable        Type = "UNAVAILABLE"         // Сервис недоступен
	TypeDataLoss           Type = "DATA_LOSS"           // Невозможно восстановить данные
	TypeCanceled           Type = "Canceled"
)

// Type тип (класс) сбоя.
type Type string

// IsVoid возвращает true, если тип - пустой.
func (t Type) IsVoid() bool {
	return t == ""
}
