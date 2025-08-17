// Package passwords для всего того, что связано с паролями: генерация, валидация, хеширование и т.д.
package passwords

import "errors"

var ErrValidationFailed = errors.New("password has invalid format")

type (
	// Encryptor предоставляет методы хеширования и верификации паролей.
	Encryptor interface {
		// Encrypt возвращает хеш для указанной строки пароля, если не удалось рассчитать хеш - вернется пустая строка и ошибка.
		Encrypt(password string) (string, error)
		// Verify вернет true, если пароль соответствует переданному хешу, ранее рассчитанному вызовом Encrypt().
		Verify(password, hash string) bool
	}

	// Validator выполняет проверку паролей на соответствие заданным правилам\ограничениям.
	Validator interface {
		// Check вернет error, производный от ErrValidationFailed, если пароль не прошел валидацию.
		Check(password string) error
	}
)
