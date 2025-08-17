package checks

import (
	"errors"
	"fmt"

	"golang.org/x/exp/constraints"
)

var (
	ErrIsNotSpecified  = errors.New("is not specified")
	ErrExceededMaximum = errors.New("exceeded maximum")
	ErrExceededMinimum = errors.New("exceeded minimum")
)

type (
	Check func() error
)

func (c Check) Check() error {
	if c == nil {
		return nil
	}

	return c()
}

func (c Check) WithAttribute(name string) Check {
	return CheckAttribute(name, c)
}

func CheckAttribute(name string, check Check) Check {
	return func() error {
		if err := check(); err != nil {
			return fmt.Errorf("%w: %s", err, name)
		}

		return nil
	}
}

func CheckJoin(checks ...Check) Check {
	return func() error {
		for _, c := range checks {
			if err := c(); err != nil {
				return err
			}
		}

		return nil
	}
}

func CheckIsDefined[T comparable](value *T) Check {
	return func() error {
		if err := CheckIsNonNil(value).Check(); err != nil {
			return err
		}

		return CheckIsNonEmpty(*value).Check()
	}
}

func CheckIsNonEmpty[T comparable](v T) Check {
	return func() error {
		var blank T
		if v == blank {
			return ErrIsNotSpecified
		}

		return nil
	}
}

func CheckIsNonNil[T any](v *T) Check {
	return func() error {
		if v == nil {
			return ErrIsNotSpecified
		}

		return nil
	}
}

func CheckMin[T constraints.Ordered](v T, n T) Check {
	return func() error {
		if v < n {
			return fmt.Errorf("%w: %v", ErrExceededMinimum, n)
		}

		return nil
	}
}

func CheckMax[T constraints.Ordered](v T, n T) Check {
	return func() error {
		if v > n {
			return fmt.Errorf("%w: %v", ErrExceededMaximum, n)
		}

		return nil
	}
}

func CheckMaxLength[T ~string](v T, n int) Check {
	return func() error {
		if len([]rune(v)) > n {
			return fmt.Errorf("%w: length %d", ErrExceededMaximum, n)
		}

		return nil
	}
}

func CheckMaxSize[T any](v []T, n int) Check {
	return func() error {
		if len(v) > n {
			return fmt.Errorf("%w: size %d", ErrExceededMaximum, n)
		}

		return nil
	}
}
