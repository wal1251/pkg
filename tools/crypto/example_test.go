package crypto

import (
	"fmt"
)

func ExampleNewHMAC() {
	// Создаем секретный ключ для HMAC.
	secret := []byte("secret")

	// Создаем экземпляр HMAC.
	hmac := NewHMAC(string(secret))

	// Задаем данные, которые мы хотим подписать.
	data := "data"

	// Подписываем данные, используя ранее созданный экземпляр HMAC.
	// Это приводит к генерации уникальной подписи для данных с использованием секретного ключа.
	signed := hmac.Sign(data)

	fmt.Println(signed)

	// Output:
	// 1b2c16b75bd2a870c114153ccda5bcfca63314bc722fa160d690de133ccbb9db
}
