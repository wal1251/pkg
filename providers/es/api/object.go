package api

import (
	"io"

	"github.com/wal1251/pkg/tools/checks"
	serial "github.com/wal1251/pkg/tools/serial"
)

type Object map[string]any

func (o Object) Object(property string) Object {
	if o == nil {
		return nil
	}

	return AsObject(o[property])
}

func (o Object) Copy() Object {
	if o == nil {
		return nil
	}

	object := make(Object)

	for key, value := range o {
		switch {
		case IsObject(value):
			object[key] = AsObject(value).Copy()
		case IsArray(value):
			object[key] = AsArray(value).Copy()
		default:
			object[key] = value
		}
	}

	return object
}

func (o Object) ApplyFrom(src Object) Object {
	if o == nil {
		return src
	}

	for key, value := range src {
		if _, ok := o[key]; ok {
			switch {
			case IsObject(value) && IsObject(o[key]):
				o[key] = AsObject(o[key]).ApplyFrom(AsObject(value))
			case IsArray(value) && IsArray(o[key]):
				o[key] = AsArray(o[key]).Merge(AsArray(value))
			default:
				o[key] = value
			}
		} else {
			o[key] = value
		}
	}

	return o
}

func (o Object) Merge(object Object) Object {
	if o == nil {
		return object
	}

	return o.Copy().ApplyFrom(object)
}

func (o Object) Array(property string) Array {
	if o == nil {
		return nil
	}

	return AsArray(o[property])
}

func (o Object) IsObject(property string) bool {
	if o == nil {
		return false
	}

	return IsObject(o[property])
}

func (o Object) IsArray(property string) bool {
	if o == nil {
		return false
	}

	return IsArray(o[property])
}

func (o Object) JSONReader() io.Reader {
	if o == nil {
		return nil
	}

	return NewJSONReader(o)
}

func (o Object) String() string {
	if o == nil {
		return ""
	}

	raw, err := serial.ToBytes(o, serial.JSONEncode[Object])
	if err != nil {
		return ""
	}

	return string(raw)
}

func AsObject(object any) Object {
	if object == nil {
		return nil
	}

	if v, ok := object.(map[string]any); ok {
		return v
	}

	if v, ok := object.(Object); ok {
		return v
	}

	return nil
}

func IsObject(object any) bool {
	if object == nil {
		return false
	}

	_, ok := object.(map[string]any)
	if !ok {
		_, ok = object.(Object)
	}

	return ok
}

func CreateObject(options ...func(Object)) Object {
	object := make(Object)
	for _, option := range options {
		option(object)
	}

	return object
}

func Property(key string, value any) func(Object) {
	return func(object Object) {
		if object != nil {
			object[key] = value
		}
	}
}

func PropertyOmitempty(key string, value any) func(Object) {
	return func(object Object) {
		if object != nil && !checks.IsEmptyAny(value) {
			object[key] = value
		}
	}
}
