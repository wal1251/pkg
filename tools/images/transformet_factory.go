package images

import (
	"fmt"
	"reflect"
)

var _ TransformerRepo = (TransformerFactory)(nil)

type (
	TransformerFactory map[string]reflect.Type
)

func (r TransformerFactory) Register(v any) {
	typ := reflect.ValueOf(v).Type()
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	transformerType := reflect.TypeOf((*Transformer)(nil)).Elem()
	if !typ.Implements(transformerType) {
		panic(fmt.Sprintf("must implement %v", transformerType))
	}

	r[typ.String()] = typ
}

func (r TransformerFactory) Create(name string) Transformer {
	if t, ok := r[name]; ok {
		if transformer, ok := reflect.New(t).Interface().(Transformer); ok {
			return transformer
		}

		return (Transformer)(nil)
	}

	return nil
}
