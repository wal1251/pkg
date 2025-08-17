package clock

import (
	"time"
)

// RetryingWrapper возвращает функцию-обёртку, которая позволяет выполнять повторные попытки
// вызова переданной функции handler.
// Параметры retries и interval определяют количество повторных попыток и интервал между ними.
// Функция onErr вызывается при каждой ошибке, возвращаемой handler, и если она возвращает false,
// повторные попытки прекращаются.
//
// Параметры:
// retries - количество повторных попыток.
// interval - интервал между попытками.
//
// Возвращаемые значения:
// Возвращает функцию, которая принимает две функции: handler (функция для повторного вызова) и
// onErr (функция, вызываемая при каждой ошибке, возвращаемой handler).
// Возвращает ошибку, если все попытки завершились неудачно.
func RetryingWrapper(retries int, interval time.Duration) func(handler func() error, onErr func(error) bool) error {
	return func(handler func() error, onErr func(error) bool) error {
		var err error
		for retry := 0; retry < retries+1; retry++ {
			if err = handler(); err == nil {
				return nil
			}

			if onErr != nil && !onErr(err) {
				break
			}

			SleepFuzz(interval * time.Duration(retry+1))
		}

		return err
	}
}

// RetryingWrapperWE возвращает функцию-обёртку, которая позволяет выполнять повторные попытки
// вызова переданной функции retryableFunc.
// Эта функция принимает количество попыток и интервал времени между попытками.
//
// Если функция (retryableFunc) не возвращает ошибку, то повторяющиеся попытки прекращаются.
// Если функция (retryableFunc) возвращает ошибку и флаг (false), указывающий на отсутствие
// необходимости повторения, то повторяющиеся попытки прекращаются.
// В других случаях повторяющиеся попытки продолжаются.
//
// Отличие от RetryingWrapper: В RetryingWrapperWE пользовательская функция retryableFunc возвращает пару значений (bool, error),
// где булевый флаг указывает, следует ли продолжать попытки. В отличие от RetryingWrapper, где повторные попытки контролируются
// внешней функцией onErr.
//
// Параметры:
// retries - количество попыток выполнения функции.
// interval - временной интервал между попытками.
//
// Возвращаемое значение:
// Возвращает функцию, которая в свою очередь принимает пользовательскую функцию retryableFunc.
// Эта внутренняя функция возвращает ошибку, если все попытки завершились неудачей или если retryableFunc вернула ошибку.
func RetryingWrapperWE(retries int, interval time.Duration) func(retryableFunc func() (bool, error)) error {
	return func(retryableFunc func() (bool, error)) error {
		var err error
		for i := 0; i < retries+1; i++ {
			if err != nil {
				_ = SleepFuzz(interval * time.Duration(i))
			}

			var retry bool
			if retry, err = retryableFunc(); err == nil {
				return nil
			}
			if !retry {
				break
			}
		}

		return err
	}
}
