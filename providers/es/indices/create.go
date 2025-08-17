package indices

import (
	"fmt"
	"reflect"

	"github.com/wal1251/pkg/providers/es/api"
)

type Create struct {
	Mappings *Mapping `json:"mappings"`
}

func IndexCreate(opts ...api.FactoryOption[Create]) api.Factory[Create] {
	return api.Factory[Create](func() (Create, error) {
		return Create{}, nil
	}).With(opts...)
}

func WithMappingsOf(t reflect.Type, opts ...api.FactoryOption[Mapping]) api.FactoryOption[Create] {
	return WithMappings(CreateMapping(OfType(t)).With(opts...))
}

func WithMappings(f api.Factory[Mapping]) api.FactoryOption[Create] {
	return func(create *Create) error {
		mapping, err := f.Produce()
		if err != nil {
			return fmt.Errorf("can't create mapping: %w", err)
		}

		create.Mappings = &mapping

		return nil
	}
}
