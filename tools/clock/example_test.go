package clock

import (
	"errors"
	"fmt"
	"time"

	"github.com/wal1251/pkg/core/errs"
)

func ExampleRetryingWrapper() {
	var counter int

	// Максимум 3 попытки, с шагом ожидания 100мс.
	retry := RetryingWrapper(3, 100*time.Millisecond)

	// Напечатаем ошибку, если она не nil, в противном случае OK.
	printErr := func(err error, s string) {
		if err == nil {
			fmt.Println(s, "OK")
		} else {
			fmt.Println(s, err)
		}
	}

	// ОК, т.к. не превысили лимит попыток.
	counter = 0
	printErr(retry(func() error {
		fmt.Println("example 1: retries count", counter)
		// Со второй повторной попытки вернем отсутствие ошибки.
		if counter == 2 {
			return nil
		}
		counter++
		return errors.New("fake")
	}, nil), "example 1:")

	// Ошибка, т.к. превысили лимит попыток.
	counter = 0
	printErr(retry(func() error {
		fmt.Println("example 2: retries count", counter)
		counter++
		return errors.New("fake")
	}, nil), "example 2:")

	// Ошибка, т.к. прервали принудительно повторные попытки.
	counter = 0
	printErr(retry(func() error {
		fmt.Println("example 3: retries count", counter)
		// Со второй повторной попытки прервем последующие попытки.
		if counter == 2 {
			return errs.ErrCancelled
		}
		counter++
		return errors.New("fake")
	}, func(err error) bool {
		// Повторять, пока не получим errs.ErrCancelled.
		return !errors.Is(err, errs.ErrCancelled)
	}), "example 3:")

	// Output:
	// example 1: retries count 0
	// example 1: retries count 1
	// example 1: retries count 2
	// example 1: OK
	// example 2: retries count 0
	// example 2: retries count 1
	// example 2: retries count 2
	// example 2: retries count 3
	// example 2: fake
	// example 3: retries count 0
	// example 3: retries count 1
	// example 3: retries count 2
	// example 3: CANCELLED
}

func ExampleRetryingWrapperWE() {
	// Функция, которая может временно терпеть неудачу.
	var count int
	retryableFunc := func() (bool, error) {
		count++
		if count < 3 {
			// возвращаем true, чтобы продолжить попытки повторения.
			return true, fmt.Errorf("error") // Повторяем
		}
		// возвращаем false, чтобы прервать попытки.
		return false, nil // Успех
	}

	// Оборачиваем функцию с 3 попытками повторения и интервалом в 100мс
	// между ними.
	wrappedFunc := RetryingWrapperWE(5, 100*time.Millisecond)

	// Вызываем обернутую функцию.
	err := wrappedFunc(retryableFunc)
	if err != nil {
		fmt.Println("Something went wrong:", err)
	} else {
		fmt.Println("Successfully executed")
	}

	// Output: Successfully executed
}
