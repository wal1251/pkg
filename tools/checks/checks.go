package checks

import (
	"reflect"
)

// IsEmpty проверяет, равно ли значение value его нулевому значению.
// Параметры:
// value T - переменная для проверки, где T - любой сравнимый тип.
//
// Возвращаемое значение:
// Возвращает true, если значение value равно его нулевому значению, иначе false.
func IsEmpty[T comparable](value T) bool {
	var blank T

	return value == blank
}

// IsEmptyAny проверяет, является ли значение value "пустым" в зависимости от его типа.
// Для различных типов определение "пустого" значения различается:
// - Для строк, массивов, слайсов и карт "пустым" считается значение с длиной 0.
// - Для булевых значений "пустым" считается false.
// - Для числовых типов (целых и с плавающей точкой) "пустым" считается значение 0.
// - Для указателей, интерфейсов "пустым" считается nil.
// Для других типов, которые не входят в эти категории, функция всегда возвращает false,
// поскольку нет универсального способа определить, является ли значение "пустым".
//
// Параметры:
// value - значение любого типа, которое необходимо проверить.
//
// Возвращаемое значение:
// Возвращает true, если значение считается "пустым" для своего типа, иначе false.
func IsEmptyAny(value any) bool {
	reflectValue := reflect.ValueOf(value)
	switch reflectValue.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return reflectValue.Len() == 0
	case reflect.Bool:
		return !reflectValue.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflectValue.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return reflectValue.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return reflectValue.Float() == 0
	case reflect.Interface, reflect.Pointer:
		return reflectValue.IsNil()
	default:
		return false
	}
}

// MustBePositive проверяет, что целочисленное значение больше 0.
// Функция вызывает панику, если значение меньше или равно 0.
//
// Параметры:
// value - значение для проверки.
// paramName - имя параметра, используется в сообщении паники.
func MustBePositive(value int, paramName string) {
	if value <= 0 {
		panic(paramName + " must be greater than 0")
	}
}
