package api

import (
	"io"

	"github.com/wal1251/pkg/tools/serial"
)

type Array []any

func (a Array) Object(index int) Object {
	if a == nil {
		return nil
	}

	return AsObject(a[index])
}

func (a Array) Array(index int) Array {
	if a == nil {
		return nil
	}

	return AsArray(a[index])
}

func (a Array) IsObject(index int) bool {
	if a == nil {
		return false
	}

	return IsObject(a[index])
}

func (a Array) IsArray(index int) bool {
	if a == nil {
		return false
	}

	return IsArray(a[index])
}

func (a Array) Copy() Array {
	if a == nil {
		return nil
	}

	result := make(Array, len(a))

	for index, value := range a {
		switch {
		case IsObject(value):
			result[index] = AsObject(value).Copy()
		case IsArray(value):
			result[index] = AsArray(value).Copy()
		default:
			result[index] = value
		}
	}

	return result
}

func (a Array) Merge(op Array) Array {
	if a == nil {
		return op
	}

	return append(a.Copy(), op...)
}

func (a Array) JSONReader() io.Reader {
	if a == nil {
		return nil
	}

	return NewJSONReader(a)
}

func (a Array) String() string {
	if a == nil {
		return ""
	}

	raw, err := serial.ToBytes(a, serial.JSONEncode[Array])
	if err != nil {
		return ""
	}

	return string(raw)
}

func AsArray(arr any) Array {
	if arr == nil {
		return nil
	}

	if v, ok := arr.([]any); ok {
		return v
	}

	if v, ok := arr.(Array); ok {
		return v
	}

	return nil
}

func IsArray(arr any) bool {
	if arr == nil {
		return false
	}

	_, ok := arr.([]any)
	if !ok {
		_, ok = arr.(Array)
	}

	return ok
}
