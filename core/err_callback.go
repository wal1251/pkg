package core

var _ ErrorCallback = (ErrorCallbackFn)(nil)

type (
	// ErrorCallback интерфейс используется для реализации обратных вызовов оповещения программных компонентов о
	// возникающих ошибках.
	ErrorCallback interface {
		// OnError используется для оповещения программных компонентов об ошибке. По соглашению, если метод возвращает
		// true - это является сигналом вызывающему, что ошибка обработана, можно продолжить выполнение, если же false -
		// вызывающему необходимо прервать выполнение операции. При получении nil в качестве ошибки метод должен вернуть
		// true.
		OnError(error) bool
	}

	// ErrorCallbackFn каноничная функция обратного вызова для оповещения о произошедшей ошибке. Реализует ErrorCallback.
	ErrorCallbackFn func(err error) bool
)

func (c ErrorCallbackFn) OnError(err error) bool {
	if c == nil {
		return err == nil
	}

	return c(err)
}

func ErrIntercept(err error, callback ErrorCallback) error {
	if err != nil && callback.OnError(err) {
		return nil
	}

	return err
}

func ErrNotify(err error, callback ErrorCallback) bool {
	if err != nil {
		return callback.OnError(err)
	}

	return true
}
