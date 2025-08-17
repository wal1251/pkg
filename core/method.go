package core

import "context"

var _ Map[any, any] = NilMap

type (
	// Method обобщенный метод, для приведения любой функции или метода приложения к каноничному виду. Принимает на вход
	// контекст вызывающего и аргумент, возвращает результат выполнения или ошибку. Если метод не имеет аргумента или
	// возвращаемого значения, тогда в качестве соответствующего параметра типа можно использовать any, а в качестве
	// значения nil.
	Method[ARG, RESULT any] func(context.Context, ARG) (RESULT, error)

	// Map обобщенная функция преобразования значения, принимает на вход значение и возвращает преобразованное значение
	// в качестве результата.
	Map[ARG, RESULT any] func(ARG) RESULT
)

// Call выполняет вызов метода, если он не равен nil. В противном случае вернет пустое значение в качестве результата.
func (m Method[ARG, RESULT]) Call(ctx context.Context, arg ARG) (RESULT, error) {
	if m == nil {
		var blank RESULT

		return blank, nil
	}

	return m(ctx, arg)
}

// Map выполнит преобразование, если функция не равна nil. В противном случае вернет пустое значение в качестве результата.
func (f Map[ARG, RESULT]) Map(arg ARG) RESULT {
	if f == nil {
		var blank RESULT

		return blank
	}

	return f(arg)
}

// MapMethodParameters предназначена для вызова Method в случае, если типы входного и\или выходного параметров отличаются
// от заданных для метода. Функция mapArg выполняет преобразование аргумента, а mapResult - результата.
func MapMethodParameters[AI, AO, RI, RO any](m Method[AI, RI], mapArg Map[AO, AI], mapResult Map[RI, RO]) Method[AO, RO] {
	return func(ctx context.Context, to AO) (RO, error) {
		result, err := m.Call(ctx, mapArg.Map(to))
		if err != nil {
			var blank RO

			return blank, err
		}

		return mapResult.Map(result), nil
	}
}

// NilMap заглушка для преобразования пустого параметра.
func NilMap(any) any {
	return nil
}
