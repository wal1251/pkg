package passwords

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var _ Encryptor = (*BCryptEncryptor)(nil)

type BCryptEncryptor struct{}

// Encrypt возвращает хеш для указанной строки пароля, если не удалось рассчитать хеш - вернется пустая строка и ошибка.
func (e *BCryptEncryptor) Encrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("can't encrypt password: %w", err)
	}

	return string(hash), nil
}

// Verify вернет true, если пароль соответствует переданному хешу, ранее рассчитанному вызовом Encrypt().
func (e *BCryptEncryptor) Verify(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// NewBcryptEncryptor возвращает новый экземпляр BCryptEncryptor.
func NewBcryptEncryptor() *BCryptEncryptor {
	return &BCryptEncryptor{}
}
