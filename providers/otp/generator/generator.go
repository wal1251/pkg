// Package generator
// Данный пакет реализует генератор OTP.
package generator

// Generator - интерфейс, который определяет контракт для генерации одноразовых паролей.
type Generator interface {
	Generate() (string, error)
}
