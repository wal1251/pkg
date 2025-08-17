package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Config хранит конфигурацию для создания HMAC подписи.
type Config struct {
	SecretSignature []byte
}

// NewHMAC инициализирует и возвращает новый экземпляр Config с заданным секретным ключом.
//
// Параметр secret представляет собой строку, которая преобразуется в байты
// и используется как ключ для HMAC.
func NewHMAC(secret string) *Config {
	return &Config{
		SecretSignature: []byte(secret),
	}
}

// Sign генерирует HMAC подпись для данного сообщения используя алгоритм SHA256.
func (c Config) Sign(message string) string {
	mac := hmac.New(sha256.New, c.SecretSignature)
	mac.Write([]byte(message))

	return hex.EncodeToString(mac.Sum(nil))
}
