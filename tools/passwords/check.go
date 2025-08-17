package passwords

import (
	"fmt"
	"regexp"
)

var _ Validator = (Check)(nil)

// Check вернет error, производный от ErrValidationFailed, если пароль не прошел валидацию.
type Check func(password string) error

// Check см. Validator.Check().
func (c Check) Check(password string) error {
	if c != nil {
		return c(password)
	}

	return nil
}

// CheckContains проверяет, что пароль содержит указанный паттерн. В противном случае вернет ошибку, производную от
// ErrValidationFailed, с уточнением ошибки msg.
func CheckContains(pattern, msg string) Check {
	regex := regexp.MustCompile(pattern)

	return func(password string) error {
		if regex.FindString(password) == "" {
			return fmt.Errorf("%w: %s", ErrValidationFailed, msg)
		}

		return nil
	}
}

// CheckNotContains проверяет, что пароль не содержит указанного паттерна. В противном случае вернет ошибку, производную
// от ErrValidationFailed, с уточнением ошибки msg.
func CheckNotContains(pattern, msg string) Check {
	regex := regexp.MustCompile(pattern)

	return func(password string) error {
		if regex.FindString(password) == "" {
			return nil
		}

		return fmt.Errorf("%w: %s", ErrValidationFailed, msg)
	}
}

// CheckLength проверяет, что длина пароля находится в заданных границах (от min до max).
func CheckLength(min, max int) Check {
	return func(password string) error {
		length := len([]rune(password))
		if length > max {
			return fmt.Errorf("%w: length exeeds %d", ErrValidationFailed, max)
		}

		if length < min {
			return fmt.Errorf("%w: length is less than %d", ErrValidationFailed, min)
		}

		return nil
	}
}

// ChecksCombine возвращает Check, который объединяет в себе заданные проверки. Применяются в порядке, как они указаны в
// аргументах.
func ChecksCombine(check Check, checks ...Check) Check {
	for _, item := range checks {
		next := item
		current := check

		check = func(password string) error {
			if err := current(password); err != nil {
				return err
			}

			return next(password)
		}
	}

	return check
}
