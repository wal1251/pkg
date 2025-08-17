package errs_test

import (
	"errors"
	"fmt"

	"github.com/wal1251/pkg/core/errs"
)

func ExampleWrapf() {
	op := func() error {
		return errs.Wrapf(errs.ErrNotImplemented, "operation is not implemented")
	}

	if err := op(); err != nil {
		if errors.Is(err, errs.ErrNotImplemented) {
			fmt.Println(err, "/Nothing to do")
			return
		}
		fmt.Println(err, "went wrong")
	}
	// Output: NOT_IMPLEMENTED: operation is not implemented /Nothing to do
}

func ExampleError_struct() {
	// Функция возвращает кастомную ошибку.
	fail := func() error {
		return errs.Error{Code: "ERR_MARKER"}
	}

	// Проверим, это та самая ошибка?
	if err := fail(); errors.Is(err, errs.Error{Code: "ERR_MARKER"}) {
		fmt.Println(err)
		return
	}

	// Сюда не должны попасть.
	fmt.Println("OK")

	// Output:
	// ERR_MARKER
}

func ExampleError_static() {
	// Функция возвращает статическую ошибку.
	fail := func() error {
		return errs.ErrForbidden
	}

	// Проверим, это та самая ошибка?
	if err := fail(); errors.Is(err, errs.ErrForbidden) {
		fmt.Println(err)
		return
	}

	// Сюда не должны попасть.
	fmt.Println("OK")

	// Output:
	// FORBIDDEN
}

func ExampleError_wrapped() {
	// Функция возвращает классифицируемую ошибку в обертке.
	fail := func() error {
		return errs.Wrapf(errs.ErrIllegalArgument, "bad request")
	}

	// Все неизвестные ошибки классифицируются как
	if err := fail(); errors.Is(err, errs.ErrIllegalArgument) {
		fmt.Println(err)
		return
	}

	// Сюда не должны попасть.
	fmt.Println("OK")

	// Output:
	// ILLEGAL_ARGUMENT: bad request
}
