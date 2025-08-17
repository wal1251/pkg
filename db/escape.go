package db

import (
	"strings"
)

// EscapePattern экранирует value символом экранирования по-умолчанию ('/'), для использования в функциях-условиях
// по паттерну. Второе возвращаемое значение означает, был ли экранирован value или нет.
func EscapePattern(value string) (string, bool) {
	var cnt int

	for i := range value {
		if c := value[i]; c == '%' || c == '_' || c == '\\' {
			cnt++
		}
	}

	// Нечего экранировать.
	if cnt == 0 {
		return value, false
	}

	var builder strings.Builder

	builder.Grow(len(value) + cnt)

	for i := range value {
		if c := value[i]; c == '%' || c == '_' || c == '\\' {
			builder.WriteByte('\\')
		}

		builder.WriteByte(value[i])
	}

	return builder.String(), true
}
