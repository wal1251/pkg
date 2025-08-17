package passwords

import "fmt"

func ExampleBCryptEncryptor() {
	// Создаем новый экземпляр BCryptEncryptor.
	var encryptor Encryptor
	encryptor = NewBcryptEncryptor()

	// Объявляем пароль.
	password := "fooBarBaz"

	// Хешируем пароль.
	encrypted, err := encryptor.Encrypt(password)
	if err != nil {
		fmt.Printf("Error occurred while encrypting password: %s\n", err)
	}

	// Проверяем соответствие переданного пароля рассчитанному хешу.
	verified := encryptor.Verify(password, encrypted)
	fmt.Println(verified)

	// Output:
	// true
}

func ExampleCheckContains() {
	var validator Validator
	validator = CheckContains(`[a-z]+`, "must contains at least one lowercase letter")

	// Объявляем список паролей для проверки.
	passwords := []string{"foobar", "fooBar", "FOO", "42"}
	for _, p := range passwords {
		// Проверяем, что пароль содержит паттерн (содержит хотя бы одну строчную букву).
		err := validator.Check(p)
		if err != nil {
			fmt.Printf("'%s' validation failed: %s\n", p, err)
			continue
		}
		fmt.Printf("'%s' passed validation\n", p)
	}

	// Output:
	// 'foobar' passed validation
	// 'fooBar' passed validation
	// 'FOO' validation failed: password has invalid format: must contains at least one lowercase letter
	// '42' validation failed: password has invalid format: must contains at least one lowercase letter
}

func ExampleCheckNotContains() {
	var validator Validator
	validator = CheckNotContains(`[a-z]+`, "must not contain lowercase letters")

	// Объявляем список паролей для проверки.
	passwords := []string{"foobar", "fooBar", "FOO", "42"}
	for _, p := range passwords {
		// Проверяем, что пароль не содержит паттерн (не содержит строчных букв).
		err := validator.Check(p)
		if err != nil {
			fmt.Printf("'%s' validation failed: %s\n", p, err)
			continue
		}
		fmt.Printf("'%s' passed validation\n", p)
	}

	// Output:
	// 'foobar' validation failed: password has invalid format: must not contain lowercase letters
	// 'fooBar' validation failed: password has invalid format: must not contain lowercase letters
	// 'FOO' passed validation
	// '42' passed validation
}

func ExampleCheckLength() {
	var validator Validator
	validator = CheckLength(4, 8)

	// Объявляем список паролей для проверки.
	passwords := []string{"foo", "fooB", "fooBar", "fooBarFo", "fooBarFooBar"}
	for _, p := range passwords {
		// Проверяем, что пароль не меньше 8 и не больше 16 символов.
		err := validator.Check(p)
		if err != nil {
			fmt.Printf("'%s' validation failed: %s\n", p, err)
			continue
		}
		fmt.Printf("'%s' passed validation\n", p)
	}

	// Output:
	// 'foo' validation failed: password has invalid format: length is less than 4
	// 'fooB' passed validation
	// 'fooBar' passed validation
	// 'fooBarFo' passed validation
	// 'fooBarFooBar' validation failed: password has invalid format: length exeeds 8
}

func ExampleChecksCombine() {
	firstValidator := CheckLength(6, 15) // Диапазон длины пароля
	validators := []Check{
		CheckContains(`\d`, "must contain at least one digit"),
		CheckContains(`[!@#$%^&*]`, "must contain at least one special character"),
		CheckContains(`[a-z]`, "must contain at least one lowercase letter"),
		CheckContains(`[A-Z]`, "must contain at least one uppercase letter"),
	}

	// Объявляем расширенный список паролей для проверки.
	passwords := []string{
		"Valid1@",           // Соответствует всем критериям
		"NoDigits!",         // Нет цифр
		"Shrt1",             // Слишком короткий
		"TooLongPassword1@", // Слишком длинный
		"NoSpecialChar1",    // Нет специальных символов
		"nouppercase1@",     // Нет заглавных букв
		"NOLOWERCASE1@",     // Нет строчных букв
	}

	for _, p := range passwords {
		// Проверяем, что пароль соответствует всем критериям.
		err := ChecksCombine(firstValidator, validators...).Check(p)
		if err != nil {
			fmt.Printf("'%s' validation failed: %s\n", p, err)
			continue
		}
		fmt.Printf("'%s' passed validation\n", p)
	}

	// Output:
	// 'Valid1@' passed validation
	// 'NoDigits!' validation failed: password has invalid format: must contain at least one digit
	// 'Shrt1' validation failed: password has invalid format: length is less than 6
	// 'TooLongPassword1@' validation failed: password has invalid format: length exeeds 15
	// 'NoSpecialChar1' validation failed: password has invalid format: must contain at least one special character
	// 'nouppercase1@' validation failed: password has invalid format: must contain at least one uppercase letter
	// 'NOLOWERCASE1@' validation failed: password has invalid format: must contain at least one lowercase letter
}
