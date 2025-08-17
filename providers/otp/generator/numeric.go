package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
)

const maxRandomDigit = 9

// NumericOTPGenerator - структура, представляющая генератор одноразовых паролей, состоящих из цифр.
type NumericOTPGenerator struct {
	length int
}

// Generate генерирует новый одноразовый пароль, состоящий из цифр.
// Возвращает сгенерированный пароль в виде строки.
func (d *NumericOTPGenerator) Generate() (string, error) {
	var otp string

	for i := 0; i < d.length; i++ {
		randomInt, err := rand.Int(rand.Reader, big.NewInt(maxRandomDigit))
		if err != nil {
			return "", fmt.Errorf("failed to generate NumericOTPGenerator: %w", err)
		}
		otp += strconv.Itoa(int(randomInt.Int64()))
	}

	return otp, nil
}

// NewNumericOTPGenerator создает новый генератор одноразовых паролей,
// состоящих из цифр указанной длины.
func NewNumericOTPGenerator(length int) *NumericOTPGenerator {
	return &NumericOTPGenerator{length}
}
