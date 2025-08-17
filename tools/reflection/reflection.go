package reflection

import (
	"reflect"
	"strings"
)

// TagParseKeyValue парсит тэги поля структуры и возвращает первое значение и словарь значений.
// "query, t1=1, t2=3" -> ("query", map{"t1": "1", "t2": "3"}).
func TagParseKeyValue(tags string) (string, map[string]string) {
	firstValue := ""
	keyValue := make(map[string]string)

	elements := strings.Split(tags, ",")
	if len(elements) == 0 {
		return firstValue, keyValue
	}

	if v := strings.TrimSpace(elements[0]); !strings.Contains(v, "=") {
		if v != "" {
			firstValue = v
		}

		elements = elements[1:]
	}

	for _, element := range elements {
		parts := strings.Split(element, "=")
		if len(parts) == 0 {
			continue
		}

		if len(parts) == 1 {
			keyValue[strings.TrimSpace(parts[0])] = ""
		}

		keyValue[strings.TrimSpace(parts[0])] = parts[1]
	}

	return firstValue, keyValue
}

// GetJSONName возвращает имя поля структуры в json формате.
func GetJSONName(field reflect.StructField) (string, bool) {
	tagValue := field.Tag.Get("json")
	if tagValue == "" {
		return field.Name, true
	}

	values := strings.Split(tagValue, ",")

	value0 := strings.TrimSpace(values[0])
	if value0 == "-" {
		return "", false
	}

	return value0, true
}

// MapByTypeName преобразует список значений в map, где key - тип значения, value - исходное значение.
// []any{1, "abc", 1.1} -> map{"int": 1, "string": "abc", "float64": 1.1}.
func MapByTypeName[T any](values ...T) map[string]T {
	dict := make(map[string]T)

	for _, value := range values {
		dict[reflect.TypeOf(value).Name()] = value
	}

	return dict
}

// IsEmptyValue The function IsEmptyValue takes in a value of any type and returns a boolean value indicating whether
// the value is empty or not. It uses reflection to determine the kind of the value and checks if it is empty based
// on its kind. For arrays, maps, slices, and strings, it checks if the length is zero. For boolean values,
// it checks if the value is false. For integer and floating-point values, it checks if the value is zero.
// For interfaces and pointers, it checks if the value is nil. If the value is of any other kind, it returns false.
func IsEmptyValue(value any) bool {
	reflectedValue := reflect.ValueOf(value)
	switch reflectedValue.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return reflectedValue.Len() == 0
	case reflect.Bool:
		return !reflectedValue.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflectedValue.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return reflectedValue.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return reflectedValue.Float() == 0
	case reflect.Interface, reflect.Pointer:
		return reflectedValue.IsNil()
	default:
		return false
	}
}
