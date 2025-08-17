package cfg

import (
	"reflect"
	"time"
)

type (
	// ValueType задает ограничение типа значения свойства конфигурации.
	ValueType interface {
		string | int | float64 | bool | time.Duration
	}

	// Property типизированное свойство конфигурации.
	Property[T ValueType] interface {
		Get() T
	}

	// PropertyProvider функция для извлечения свойства конфигурации по строке-идентификатору, может служить оберткой
	// над конкретной реализацией загрузки конфигурации.
	PropertyProvider[T ValueType] func(string) T
)

// TypeMatches вернет true, если тип переданного значения соответствует типу возвращаемого значения функции.
func (p PropertyProvider[T]) TypeMatches(v any) bool {
	if p == nil {
		return false
	}

	return reflect.TypeOf(p).Out(0).AssignableTo(reflect.TypeOf(v))
}

// Get вернет значение свойства конфигурации с помощью заданной функции.
func (p PropertyProvider[T]) Get(key string) T {
	if p == nil {
		var blank T

		return blank
	}

	return p(key)
}

// PropertyProviderAdapter функция переходник для преобразования функции PropertyProvider заданного типа к функции
// с обобщенным типом. Контроль на совместимость возвращаемых значений функций не выполняется.
func PropertyProviderAdapter[T, R ValueType](provider PropertyProvider[T]) PropertyProvider[R] {
	return func(key string) R {
		valuePtr := new(R)
		result := provider(key)
		ptr := reflect.ValueOf(valuePtr).Elem()
		ptr.Set(reflect.ValueOf(result))

		return *valuePtr
	}
}
