package gen

import (
	"sync"

	"github.com/google/uuid"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/tools/singleton"
)

// UUIDGenerator генератор случайных уникальных идентификаторов.
type UUIDGenerator func() uuid.UUID

// UUID вернет генератор уникальных идентификаторов в соответствие со средой исполнения приложения, для выполнения
// тестов будет создан моковый генератор.
func UUID() UUIDGenerator {
	if cfg.IsTestRuntime() {
		return MockUUIDGenerator()
	}

	return RealUUIDGenerator()
}

// Next генерирует очередной уникальный идентификатор.
func (g UUIDGenerator) Next() uuid.UUID {
	if g == nil {
		return uuid.Nil
	}

	return g()
}

// nolint используем синглтон для мока uuid генератора
var uuidSingleton = singleton.NewSingleton(DummyRand)

// MockUUIDGenerator возвращает моковый генератор уникальных идентификаторов для тестов.
func MockUUIDGenerator() UUIDGenerator {
	var lock sync.Mutex

	return func() uuid.UUID {
		lock.Lock()
		defer lock.Unlock()

		var result uuid.UUID
		if _, err := uuidSingleton.Get().Read(result[:]); err != nil {
			result = uuid.Nil
		}

		return result
	}
}

// RealUUIDGenerator возвращает генератор уникальных идентификаторов.
func RealUUIDGenerator() UUIDGenerator {
	return uuid.New
}
