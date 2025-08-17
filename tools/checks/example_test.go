package checks

import (
	"fmt"
)

func ExampleIsEmpty() {
	// Пример использования с различными типами.
	fmt.Println(IsEmpty(0))           // Для int
	fmt.Println(IsEmpty(""))          // Для string
	fmt.Println(IsEmpty(false))       // Для bool
	fmt.Println(IsEmpty([...]int{0})) // Для массива

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleIsEmptyAny() {
	// Проверим различные типы данных.
	fmt.Println("Empty string:", IsEmptyAny(""))
	fmt.Println("String 'Hello':", IsEmptyAny("Hello"))
	fmt.Println("Empty slice:", IsEmptyAny([]int{}))
	fmt.Println("Slice [1, 2, 3]:", IsEmptyAny([]int{1, 2, 3}))
	fmt.Println("Empty map:", IsEmptyAny(map[string]int{}))
	fmt.Println("Map with elements:", IsEmptyAny(map[string]int{"a": 1}))
	fmt.Println("Bool false:", IsEmptyAny(false))
	fmt.Println("Bool true:", IsEmptyAny(true))
	fmt.Println("Int 0:", IsEmptyAny(0))
	fmt.Println("Int 42:", IsEmptyAny(42))
	fmt.Println("Empty pointer:", IsEmptyAny((*int)(nil)))
	fmt.Println("Not empty pointer:", IsEmptyAny(new(int)))

	// Output:
	// Empty string: true
	// String 'Hello': false
	// Empty slice: true
	// Slice [1, 2, 3]: false
	// Empty map: true
	// Map with elements: false
	// Bool false: true
	// Bool true: false
	// Int 0: true
	// Int 42: false
	// Empty pointer: true
	// Not empty pointer: false
}

func ExampleMustBePositive() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in ExampleMustBePositive:", r)
		}
	}()

	// Правильное использование - не вызывает панику.
	MustBePositive(10, "quantity")

	// Неправильное использование - вызывает панику.
	MustBePositive(-5, "quantity")

	// Output:
	// Recovered in ExampleMustBePositive: quantity must be greater than 0
}
