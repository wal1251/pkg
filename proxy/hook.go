package proxy

import (
	"context"
	"runtime/debug"
)

type (
	// Hook используется как колбэк для перехвата вызовов методов.
	Hook func(ctx context.Context, object any, method string, args []any) context.Context

	// PanicHook используется как колбэк для перехвата паники в методах.
	PanicHook func(msg any, stack []byte, object any, method string, args []any) any
)

// Hook выполняет вызов h с указанными параметрами, если h не nil.
func (h Hook) Hook(ctx context.Context, object any, method string, args []any) context.Context {
	if h != nil {
		return h(ctx, object, method, args)
	}

	return ctx
}

// And добавляет цепочку последовательных вызовов hooks после вызова h.
func (h Hook) And(hooks ...Hook) Hook {
	current := h
	for _, hook := range hooks {
		c := current
		next := hook
		current = func(ctx context.Context, object any, method string, args []any) context.Context {
			return next(c(ctx, object, method, args), object, method, args)
		}
	}

	return current
}

// Hook выполняет вызов h с указанными параметрами, если h не nil.
func (h PanicHook) Hook(object any, method string, args []any) {
	if h == nil {
		return
	}

	if err := recover(); err != nil {
		if err := h(err, debug.Stack(), object, method, args); err != nil {
			panic(err)
		}
	}
}

// And добавляет цепочку последовательных вызовов panicHooks после вызова h.
func (h PanicHook) And(hooks ...PanicHook) PanicHook {
	current := h
	for _, hook := range hooks {
		c := current
		next := hook
		current = func(msg any, stack []byte, object any, method string, args []any) any {
			// Вызываем текущий обработчик паники и получаем его результат.
			res := c(msg, stack, object, method, args)
			// Передаём результат в следующий обработчик в цепочке.
			return next(res, stack, object, method, args)
		}
	}

	return current
}

// ExtractContext возвращает контекст из первого аргумента, если первый элемент является контектом, в противном случае,
// вернет context.Background. Если контекст был извлечен из первого аргумента. В качестве второго возвращаемого параметра
// вернет аргументы без извлеченного контекста.
func ExtractContext(args []any) context.Context {
	if len(args) > 0 {
		if ctx, ok := args[0].(context.Context); ok {
			return ctx
		}
	}

	return context.Background()
}

// HasError проверяет, является ли последний элемент типом error.
func HasError(results []any) bool {
	if len(results) > 0 {
		if last := results[len(results)-1]; last != nil {
			_, ok := last.(error)

			return ok
		}
	}

	return false
}

// ExtractErr извлекает error из последнего results, вторым параметром возвращает results без извлеченного error.
func ExtractErr(results []any) ([]any, error) {
	if err, ok := results[len(results)-1].(error); ok {
		return results[:len(results)-1], err
	}

	return results, nil
}
