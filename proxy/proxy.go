package proxy

import (
	"context"
	"reflect"
	"runtime/debug"
)

var _ PanicInterceptor = DefaultPanicInterceptor

type (
	// GenericFunction обобщенная функция.
	GenericFunction func([]reflect.Value) []reflect.Value

	// MethodInvocationInterceptor перехватчик вызова метода.
	MethodInvocationInterceptor func(object any, name string, args []reflect.Value) []reflect.Value

	// MethodInvocationMiddleware прокси мидлварь вызова метода.
	MethodInvocationMiddleware func(object any, name string, args []reflect.Value, next GenericFunction) []reflect.Value

	// PanicInterceptor перехватчик паники.
	PanicInterceptor func(err any, stack []byte) any

	// InvocationArguments аргументы вызова метода.
	InvocationArguments []reflect.Value

	// InvocationResults возвращенные параметры метода.
	InvocationResults []reflect.Value
)

// DefaultPanicInterceptor перехватчик паники по умолчанию.
func DefaultPanicInterceptor(err any, _ []byte) any {
	return err
}

// GenericFunction преобразование к GenericFunction со статичными аргументами object и name.
func (interceptor MethodInvocationInterceptor) GenericFunction(object any, name string) GenericFunction {
	return func(args []reflect.Value) []reflect.Value {
		return interceptor(object, name, args)
	}
}

// Context получение контекста из InvocationArguments.
func (a InvocationArguments) Context() context.Context {
	if len(a) > 0 {
		if a[0].Type().Implements(Type((*context.Context)(nil))) {
			return As[context.Context](a[0])
		}
	}

	return nil
}

// HasError возвращет true, если InvocationResults содержит error последним параметром.
func (r InvocationResults) HasError() bool {
	return len(r) > 0 && r[len(r)-1].Type().Implements(Type((*error)(nil)))
}

// GetError получает error из последнего параметра InvocationResults.
func (r InvocationResults) GetError() error {
	return As[error](r[len(r)-1])
}

// MakeMethodInvocationInterceptor возвращает MethodInvocationInterceptor содержащий последовательную цепочку вызовов middlewares.
func MakeMethodInvocationInterceptor(middlewares ...MethodInvocationMiddleware) MethodInvocationInterceptor {
	interceptor := func(object any, name string, args []reflect.Value) []reflect.Value {
		return reflect.ValueOf(object).MethodByName(name).Call(args)
	}

	for i := len(middlewares); i > 0; i-- {
		next := interceptor
		middleware := middlewares[i-1]
		interceptor = func(object any, name string, args []reflect.Value) []reflect.Value {
			return middleware(object, name, args, func(args []reflect.Value) []reflect.Value {
				return next(object, name, args)
			})
		}
	}

	return interceptor
}

// MakeFunc возвращет созданную функцию заданного типа из GenericFunction.
func MakeFunc[T any](function GenericFunction) T {
	var proxyFn T

	fn := reflect.MakeFunc(reflect.TypeOf(proxyFn), function)
	reflect.ValueOf(&proxyFn).Elem().Set(fn)

	return proxyFn
}

// As извлекает значение reflect.Value заданного типа.
func As[T any](value reflect.Value) T {
	var blank T

	if typedValue, ok := value.Interface().(T); ok {
		return typedValue
	}

	return blank
}

// Type получает значение интерфейсного типа, например:
//
//	t := Type((*context.Context)(nil)).
func Type(value any) reflect.Type {
	return reflect.TypeOf(value).Elem()
}

// PanicInterception перехват паники.
func PanicInterception(interceptor PanicInterceptor) {
	if err := recover(); err != nil {
		if interceptor != nil {
			err = interceptor(err, debug.Stack())
		}

		if err != nil {
			panic(err)
		}
	}
}
