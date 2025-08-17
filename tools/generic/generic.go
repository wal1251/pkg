package generic

// IndirectGet возвращает значение разыменованного указателя ptr, если он не nil, в противном случае вернет нулевое
// значение типа T.
func IndirectGet[T any](ptr *T) T {
	if ptr == nil {
		var z T

		return z
	}

	return *ptr
}

// IndirectGetWithDefault возвращает значение разыменованного указателя ptr, если он не nil, в противном случае вернет
// нулевое дефолтное значение deft.
func IndirectGetWithDefault[T any](ptr *T, deft T) T {
	if ptr == nil {
		return deft
	}

	return *ptr
}

// ApplyIfNotNil Применить действие с value, если он не nil.
func ApplyIfNotNil[R any, T any](value *T, f func(T) R) {
	if value != nil {
		f(*value)
	}
}

func ApplyNillable[R any, T comparable](value *T, setNillable func(*T) R, clear func() R) {
	if value != nil {
		var z T
		if *value == z {
			clear()
		} else {
			setNillable(value)
		}
	}
}

// EqualValues The function "EqualValues" takes two pointers of type T and returns a boolean value.
// T must be a comparable type. The function checks if both pointers are nil and returns true if they are.
// If only one pointer is nil, the function returns false. If both pointers have values,
// the function compares the values and returns true if they are equal, otherwise it returns false.
// This function can be used to check if two values of any comparable type are equal.
func EqualValues[T comparable](val1 *T, val2 *T) bool {
	if val1 == nil && val2 == nil {
		return true
	}
	if val1 == nil || val2 == nil {
		return false
	}

	return *val1 == *val2
}

// ParamsRetryWithBackoff параметры для функции RetryWithBackoff.
type ParamsRetryWithBackoff[T any] struct {
	// MaxRetries максимальное количество попыток повтора.
	MaxRetries int

	// ErrMapFunc функция для обработки (или маппинга) ошибки, возвращённой операцией.
	// Возвращает обработанную ошибку. Если возвращаемая ошибка nil, попытки повтора прекращаются.
	ErrMapFunc func(err error) error

	// CheckRetryNecessity функция для проверки, требуется ли повтор.
	// Если возвращает false, повторы прекращаются.
	CheckRetryNecessity func(err error) bool

	// Operation выполняемая операция, которая может быть повторена.
	// Возвращает результат и ошибку.
	Operation func() (*T, error)
}

// RetryWithBackoff выполняет операцию с заданным числом повторов. Если операция завершилась успешно, возвращает результат.
func RetryWithBackoff[T any](params ParamsRetryWithBackoff[T]) (*T, error) {
	var result *T
	var err error

	for attempt := 1; attempt <= params.MaxRetries; attempt++ {
		result, err = params.Operation()
		if err == nil {
			return result, nil
		}

		mappedErr := params.ErrMapFunc(err)
		if mappedErr == nil {
			return result, nil
		}

		if !params.CheckRetryNecessity(mappedErr) {
			return nil, mappedErr
		}
	}

	return nil, err
}
