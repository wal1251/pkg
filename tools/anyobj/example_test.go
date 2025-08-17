package anyobj

import (
	"fmt"
)

type User struct {
	Email  string
	Name   *string
	active bool
}

func (u User) String() string {
	nameStr := "nil"
	if u.Name != nil {
		nameStr = *u.Name
	}
	return fmt.Sprintf("Email: %s, Active: %t, Name: %s", u.Email, u.active, nameStr)
}

func ExampleSafeCopy() {
	// Создаем экземпляр структуры, из которой будем копировать.
	srcName := "John"
	src := User{
		Email:  "john@example.com",
		Name:   &srcName,
		active: true,
	}

	// Создаем экземпляр структуры, в которую будем копировать.
	var dest User

	// Вызываем метод безопасного копирования.
	err := SafeCopy(src, &dest)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Выведем в консоль исходную и скопированную структуру.
	fmt.Printf("Source: %s\n", src)
	fmt.Printf("Destination: %s\n", dest)

	// Проверим, что поля Name в исходной и целевой структурах не используют один и тот же указатель на строку.
	fmt.Printf("Source and destination pointer fields addresses not equal: %v\n", src.Name != dest.Name)

	// Внимание! Поля active не совпадают, т.к. это поле является неэкспортируемым и не копируется.

	// Output:
	// Source: Email: john@example.com, Active: true, Name: John
	// Destination: Email: john@example.com, Active: false, Name: John
	// Source and destination pointer fields addresses not equal: true
}
