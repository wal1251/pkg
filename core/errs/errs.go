// Package errs для работы с ошибками, классификация и перехват причин ошибок.
//
// При классификации ошибок, связанных с логикой крайне рекомендуется упаковывать ошибки с помощью функций пакета.
// Например:
//
//	err := someSearch()
//	if err != nil {
//		return errs.Wrapf(errs.ErrNotFound, err)
//	}
//
// Для того, чтобы классифицированные ошибки могли быть перехвачены в любом слое приложения.
//
// Неупакованными можно оставлять системные ошибки, например сбой ввода, вывода, ошибка соединения и т.д...
package errs

import (
	"errors"
	"fmt"
)

var _ error = (*Error)(nil)

// Error служит для классификации возникающих ошибок. Поддерживает интерфейс error. В месте возникновения нужно создать
// структуру ошибки с нужным кодом и типом. Далее эту ошибку можно перехватить с помощью стандартной библиотеки работы
// с ошибками: errors.Is() и errors.As().
//
// Например:
//
//	err := errs.Error{Code: "MARK"} // такую ошибку можно вернуть из функции.
//
//	if errors.Is(err, errs.Error{Code: "MARK"}) {
//		// обрабатываем эту ошибку ...
//	}
//
// Для создания необходимой ошибки, крайне рекомендуется пользоваться функциями хелперами или статическими ошибками
// пакета errs или определенными в приложении.
type Error struct {
	Code    string
	Type    Type
	ErrNum  string
	Details map[string]string
}

func (r Error) Error() string {
	if r.Code == "" {
		return string(r.Type)
	}

	return r.Code
}

func (r Error) WithDetails(details map[string]string) Error {
	r.Details = details

	return r
}

// Is проверяет, что аргумент типа error является эквивалентной ошибкой.
// Например:
//
//	err := someSearch() // может вернуть errs.ErrNotFound
//	if err != nil {
//		if errs.ErrNotFound.Is(err) {
//			// перехватываем ошибку ...
//		}
//		// какая то другая ошибка ...
//	}
//
// .
func (r Error) Is(err error) bool {
	var target Error
	if errors.As(err, &target) {
		return r.Equals(target)
	}

	var targetPtr *Error
	if errors.As(err, &targetPtr) {
		return targetPtr != nil && r.Equals(*targetPtr)
	}

	return false
}

// Equals проверяет, что описания ошибок равнозначны.
func (r Error) Equals(v Error) bool {
	if r.Code != "" {
		return r.Code == v.Code
	}

	return r.Type == v.Type
}

// AsReason возвращает Error если ошибка ранее классифицировалась. в противном случае вернет ErrSystemFailure.
func AsReason(err error) Error {
	var r Error
	if errors.As(err, &r) {
		return r
	}

	return ErrSystemFailure
}

// Проверяет наличие цифровой ошибки в строке

// Reasons создает Error с указанным кодом и типом.
//
//	const (
//		UserInfoCtxExtract = "USER_INFO_NOT_FOUND_IN_CONTEXT"
//	)
//
//	var ErrUserInfoNotFoundInCtx = Reasons(UserInfoCtxExtract, TypeIllegalArgument)

// Reasons Первая цифра берется из errNum, если она не пустая, для обратной совместимости.
func Reasons(code string, t Type, errNum ...interface{}) Error {
	var err Error
	err.Code = code

	if !t.IsVoid() {
		err.Type = t
	}
	if len(errNum) > 0 {
		err.ErrNum = fmt.Sprintf("%v", errNum[0])
	}
	if err.ErrNum == "" {
		err.ErrNum = "0"
	}

	return err
}
