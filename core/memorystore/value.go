package memorystore

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wal1251/pkg/tools/serial"
)

// Value представляет собой обертку для значений получаемых из MemoryStore.
type Value struct {
	val []byte
}

// NewValue создает новый экземпляр Value.
func NewValue(val []byte) *Value {
	return &Value{val: val}
}

// RawBytes возвращает исходное значение в виде среза байт.
func (v *Value) RawBytes() []byte {
	return v.val
}

// Bytes возвращает значение в виде среза байт.
// Предварительно используя десериализацию из JSON.
func (v *Value) Bytes() ([]byte, error) {
	val, err := v.String()

	return []byte(val), err
}

// String возвращает строковое представление значения.
func (v *Value) String() (string, error) {
	val, err := serial.FromBytes(v.val, serial.JSONDecode[any])

	return fmt.Sprint(val), err
}

// Int преобразует значение в целое число.
func (v *Value) Int() (int, error) {
	return serial.FromBytes(v.val, serial.JSONDecode[int])
}

// Int64 преобразует значение в целое число типа int64.
func (v *Value) Int64() (int64, error) {
	return serial.FromBytes(v.val, serial.JSONDecode[int64])
}

// Uint64 преобразует значение в беззнаковое целое число типа uint64.
func (v *Value) Uint64() (uint64, error) {
	return serial.FromBytes(v.val, serial.JSONDecode[uint64])
}

// Float32 преобразует значение в число с плавающей точкой типа float32.
func (v *Value) Float32() (float32, error) {
	return serial.FromBytes(v.val, serial.JSONDecode[float32])
}

// Float64 преобразует значение в число с плавающей точкой типа float64.
func (v *Value) Float64() (float64, error) {
	return serial.FromBytes(v.val, serial.JSONDecode[float64])
}

// Bool преобразует значение в булевый тип.
func (v *Value) Bool() (bool, error) {
	return serial.FromBytes(v.val, serial.JSONDecode[bool])
}

// Time преобразует значение в тип time.Time.
func (v *Value) Time() (time.Time, error) {
	return serial.FromBytes(v.val, serial.JSONDecode[time.Time])
}

// Struct преобразует значение в структуру, используя JSON-декодирование.
func (v *Value) Struct(dst any) error {
	err := json.Unmarshal(v.val, dst)
	if err != nil {
		return fmt.Errorf("can't unmarshal value: %w", err)
	}

	return nil
}
