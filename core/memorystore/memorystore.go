package memorystore

import (
	"context"
	"errors"
	"time"
)

const DefaultMaxBulkRequestSize = 100

var (
	ErrKeyNotFound         = errors.New("no such key")
	ErrBulkRequestTooLarge = errors.New("bulk request too large")
)

type (
	// Closer предоставляет интерфейс для компонентов, которые могут быть явно закрыты.
	// Это позволяет освободить занятые ресурсы и корректно завершить работу компонента.
	Closer interface {
		// Close закрывает компонент и освобождает все занятые им ресурсы. После вызова Close, компонент
		// не должен использоваться.
		Close(ctx context.Context)
	}

	// MemoryStore представляет абстракцию хранилища данных в памяти,
	// поддерживающего операции чтения, записи и удаления.
	MemoryStore interface {
		// Set устанавливает значение для указанного ключа с опциональной длительностью жизни.
		// Если ключ уже существует, его значение будет перезаписано.
		Set(ctx context.Context, key string, value any, expiration time.Duration) error

		// Get возвращает значение, ассоциированное с указанным ключом.
		// Если ключ не найден, возвращается ошибка ErrKeyNotFound.
		Get(ctx context.Context, key string) (*Value, error)

		// GetList возвращает значения для списка ключей. Результат представляет собой список указателей на значения.
		// Для несуществующих ключей в соответствующих позициях будет nil.
		// Если количество ключей больше чем MaxBulkRequestSize, возвращается ошибка ErrBulkRequestTooLarge.
		GetList(ctx context.Context, keys ...string) ([]*Value, error)

		// Delete удаляет указанные ключи из хранилища.
		// Возвращает количество успешно удаленных ключей.
		// Если количество ключей больше чем MaxBulkRequestSize, возвращается ошибка ErrBulkRequestTooLarge.
		Delete(ctx context.Context, keys ...string) (int, error)
	}

	// Manager расширяет интерфейс MemoryStore, добавляя к нему возможность закрытия хранилища.
	Manager interface {
		MemoryStore
		Closer
	}
)
