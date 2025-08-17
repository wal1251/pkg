// Package ctxs предоставляет базовые функции для работы с контекстом приложения.
//
// Ниже представлен пример, тога как можно запустить приложение с инициализированным контекстом.
//
//	ctx := Create(CancelOnTermSignal(nil)) // создадим контекст приложения.
//	applicationStart(ctx) // запустим приложение.
//	<-ctx.Done() // ждем сигнала завершения приложения, например по Ctrl+C.
//
// Предпочтительно работать со значениями контекста через функции:
//
//	ctx = ValuePut(ctx, "LuckySeven", 777) // новый контекст со значением.
//	luckySeven := ValueGet[int](ctx, "LuckySeven") // получим значение из контекста, ключ и тип должны строго совпадать.
//
// .
package ctxs

import (
	"context"
	"fmt"
	"os"

	"github.com/wal1251/pkg/tools/sys"
)

// KeyPrefix префикс ключей контекста приложения.
const KeyPrefix = "RMR-CORE"

type (
	// Key ключ контекста приложения.
	Key string

	// Option опция модификации контекста приложения.
	Option func(context.Context) context.Context
)

// Create создает и возвращает новый контекст приложения (от context.Background), применяет к нему указанные опции.
func Create(opts ...Option) context.Context {
	return WithOptions(context.Background(), opts...)
}

// WithOptions возвращает контекст, наследованный от указанного с примененными опциями контекста приложения.
func WithOptions(ctx context.Context, opts ...Option) context.Context {
	for _, opt := range opts {
		ctx = opt(ctx)
	}

	return ctx
}

// ValuePut возвращает новую копию контекста с установленным значением по ключу.
func ValuePut[T any](ctx context.Context, key string, value T) context.Context {
	return context.WithValue(ctx, MakeKey(key), value)
}

// ValueGet извлекает и возвращает значение заданного типа по ключу из контекста.
func ValueGet[T any](ctx context.Context, key string) T {
	value, _ := ValueCheckAndGet[T](ctx, key)

	return value
}

// ValueCheckAndGet извлекает и возвращает значение по ключу из контекста и признак наличия ключа в контексте, если
// запрашиваемого ключа в контексте нет, тогда вернет false в качестве признака.
func ValueCheckAndGet[T any](ctx context.Context, key string) (T, bool) {
	var blank T

	value := ctx.Value(MakeKey(key))
	if value == nil {
		return blank, false
	}

	if typed, ok := value.(T); ok {
		return typed, true
	}

	return blank, false
}

// MakeKey возвращает новый ключ для контекста приложения.
func MakeKey(key string) Key {
	return Key(fmt.Sprintf("%s-%s", KeyPrefix, key))
}

// CancelOnTermSignal возвращает опцию контекста приложения, отменяющую контекст в случае получения приложением
// завершающего сигнала от ОС.
func CancelOnTermSignal(hook func(os.Signal)) Option {
	return func(parent context.Context) context.Context {
		ctx, cancel := context.WithCancel(parent)

		go func() {
			defer cancel()

			sig := <-sys.ShutdownSignal()
			if hook != nil {
				hook(sig)
			}
		}()

		return ctx
	}
}
