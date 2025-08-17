package errs

// Описания общесистемных ошибок.
var (
	ErrIllegalArgument    = General(TypeIllegalArgument)    // Входной параметр или запрос не может быть принят.
	ErrAuthFailure        = General(TypeAuthFailure)        // Ошибка аутентификации.
	ErrForbidden          = General(TypeForbidden)          // Выполнение операции запрещено.
	ErrNotFound           = General(TypeNotFound)           // Объект не найден.
	ErrConflict           = General(TypeConflict)           // Операция не согласуется с текущим состоянием.
	ErrCancelled          = General(TypeCancelled)          // Операция была отменена.
	ErrHasReferences      = General(TypeHasReferences)      // Имеются ссылки на объект.
	ErrNotImplemented     = General(TypeNotImplemented)     // Отсутствует реализация операции.
	ErrSystemFailure      = General(TypeSystemFailure)      // Не классифицируемая системная ошибка.
	ErrServiceUnavailable = General(TypeServiceUnavailable) // Сервис недоступен.
	ErrTooManyRequests    = General(TypeTooManyRequests)    // Превышено количество запросов.

	ErrUnauthenticated    = General(TypeUnauthenticated)    // Не авторизован
	ErrPermissionDenied   = General(TypePermissionDenied)   // Недостаточно прав для выполнения операции
	ErrResourceExhausted  = General(TypeResourceExhausted)  // Ресурс был исчерпан
	ErrFailedPrecondition = General(TypeFailedPrecondition) // Операция не может быть выполнена
	ErrAborted            = General(TypeAborted)            // Операция была прервана
	ErrOutOfRange         = General(TypeOutOfRange)         // Операция не может быть выполнена
	ErrUnimplemented      = General(TypeUnimplemented)      // Операция не реализована
	ErrInternal           = General(TypeInternal)           // Внутренняя ошибка
	ErrUnavailable        = General(TypeUnavailable)        // Сервис недоступен
	ErrDataLoss           = General(TypeDataLoss)           // Невозможно восстановить данные
)

// General возвращает общее описание причины сбоя заданного типа t.
func General(t Type) Error {
	return Error{Code: string(t), Type: t}
}
