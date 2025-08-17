package singleton

import "fmt"

func ExampleSingleton() {
	initFunc := func() *int {
		num := 10
		return &num
	}

	// Создание экземпляра Singleton.
	s := NewSingleton(initFunc)

	// Получение значения.
	instance := s.Get()
	fmt.Printf("Value: %d\n", *instance)

	// Output:
	// Value: 10
}
