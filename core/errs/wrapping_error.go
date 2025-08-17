package errs

import "fmt"

var _ error = (*WrappingError)(nil)

// WrappingError обертка-холдер для описания ошибки. Структура поддерживает работу с интерфейсом error, а так же к ней
// применимы функции стандартной библиотеки errors.Is() и errors.As().
type WrappingError struct {
	Err     error  // Ошибка-причина.
	Message string // Человеко-читаемое описание ошибки.
	Fields  map[string]string
}

func (e *WrappingError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Err.Error(), e.Message)
	}

	return e.Message
}

// Unwrap вернет исходную ошибку-причину.
func (e *WrappingError) Unwrap() error {
	return e.Err
}

// With возвращает error с конкретизированной причиной reason.
func With(reason error, err error) *WrappingError {
	return &WrappingError{
		Err:     reason,
		Message: err.Error(),
	}
}

// Wrapf возвращает error с конкретизированной причиной reason и уточняющим сообщением.
func Wrapf(reason error, message string, args ...any) *WrappingError {
	return &WrappingError{
		Err:     reason,
		Message: fmt.Sprintf(message, args...),
	}
}

func WrapFields(reason error, message string, field string) *WrappingError {
	return &WrappingError{
		Err:    reason,
		Fields: map[string]string{field: message},
	}
}
