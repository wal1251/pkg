package cfg

import (
	"fmt"
	"strings"

	"github.com/wal1251/pkg/core"
)

var _ KeyMap = KeyIdentity

type (
	Key    string             // Идентификатор свойства конфигурации приложения.
	KeyMap core.Map[Key, Key] // Выполняет преобразование ключа по заданному правилу (например, добавить префикс и т.д.).
)

// Map возвращает новый ключ, после последовательного применения указанных функций преобразования ключа.
func (k Key) Map(mapping ...KeyMap) Key {
	key := k
	for _, m := range mapping {
		key = m(key)
	}

	return key
}

func (k Key) String() string {
	return string(k)
}

// KeyWithPrefix возвращает функцию преобразования ключа, которая добавляет к ключу префикс.
func KeyWithPrefix(prefix string) KeyMap {
	if prefix == "" {
		return KeyIdentity
	}

	prefixWithDelimiter := prefix + "_"

	return func(key Key) Key {
		if strings.HasPrefix(string(key), prefixWithDelimiter) {
			return key
		}

		return Key(fmt.Sprintf("%s%s", prefixWithDelimiter, key))
	}
}

// KeyWithSuffix возвращает функцию преобразования ключа, которая добавляет к ключу постфикс.
func KeyWithSuffix(suffix string) KeyMap {
	if suffix == "" {
		return KeyIdentity
	}

	suffixWithDelimiter := "_" + suffix

	return func(key Key) Key {
		if strings.HasSuffix(string(key), suffixWithDelimiter) {
			return key
		}

		return Key(fmt.Sprintf("%s%s", key, suffixWithDelimiter))
	}
}

// KeyIdentity возвращает точно такой же ключ, который был передан в функцию. Используется как функция-заглушка.
func KeyIdentity(key Key) Key {
	return key
}
