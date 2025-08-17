package indices

import (
	"errors"
	"fmt"

	"github.com/wal1251/pkg/providers/es/api"
)

const (
	DynamicTrue    Dynamic = "true"
	DynamicRuntime Dynamic = "runtime"
	DynamicFalse   Dynamic = "false"
	DynamicStrict  Dynamic = "strict"
)

var ErrInvalidDynamicParam = errors.New("invalid dynamic parameter")

type Dynamic string

func (d Dynamic) IsValid() bool {
	return d == DynamicTrue || d == DynamicRuntime || d == DynamicFalse || d == DynamicStrict
}

func WithDynamic(dynamic Dynamic) api.FactoryOption[Mapping] {
	return func(mapping *Mapping) error {
		if !dynamic.IsValid() {
			return fmt.Errorf("%w: %s", ErrInvalidDynamicParam, dynamic)
		}

		mapping.Dynamic = &dynamic

		return nil
	}
}
