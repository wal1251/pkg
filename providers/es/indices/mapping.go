package indices

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/wal1251/pkg/providers/es/api"
	"github.com/wal1251/pkg/tools/collections"
	"github.com/wal1251/pkg/tools/reflection"
)

var (
	ErrUnknownTagKey        = errors.New("unknown tag key")
	ErrInvalidMappingObject = errors.New("invalid mapping object")
)

type (
	Mapping struct {
		Properties Properties `json:"properties,omitempty"`
		Dynamic    *Dynamic   `json:"dynamic,omitempty"`
	}

	Properties map[string]MappingType

	MappingType struct {
		Properties Properties `json:"properties,omitempty"`
		Type       *Type      `json:"type,omitempty"`
		Analyzer   *string    `json:"analyzer,omitempty"`
	}
)

func (m *MappingType) SetType(t Type) {
	m.Type = &t
}

func (m *MappingType) LoadFromTag(tagValue string) error {
	if tagValue == "" {
		return nil
	}

	first, pairs := reflection.TagParseKeyValue(tagValue)
	if first != "" {
		typ, err := parsePropertyType(first)
		if err != nil {
			return fmt.Errorf("failed to parse es type '%s' from tag %s: %w", first, tagValue, err)
		}

		m.Type = typ
	}

	for key, value := range pairs {
		switch key {
		case "type":
			typ, err := parsePropertyType(value)
			if err != nil {
				return fmt.Errorf("failed to parse es type '%s' from tag %s: %w", value, tagValue, err)
			}

			m.Type = typ
		case "analyzer":
			analyzer := value
			m.Analyzer = &analyzer
		default:
			return fmt.Errorf("%w: unknown key %s: %s", ErrUnknownTagKey, key, tagValue)
		}
	}

	return nil
}

func CreateMapping(opts ...api.FactoryOption[Mapping]) api.Factory[Mapping] {
	return api.Factory[Mapping](func() (Mapping, error) {
		return Mapping{}, nil
	}).With(opts...)
}

func OfType(t reflect.Type) api.FactoryOption[Mapping] {
	return func(mapping *Mapping) error {
		properties := make(Properties)
		if err := detectStructProperties(properties, t, collections.NewSet[string]()); err != nil {
			return err
		}

		mapping.Properties = properties

		return nil
	}
}

func FieldMapping(field reflect.StructField) (MappingType, string, error) {
	var mappingType MappingType

	property, err := detectStructFieldMapping(&mappingType, field, collections.NewSet[string]())
	if err != nil {
		return MappingType{}, "", err
	}

	return mappingType, property, nil
}

func detectMappingType(mappingType *MappingType, typ reflect.Type, check collections.Set[string]) error { //nolint: cyclop
	switch typ.Kind() {
	case reflect.Struct:
		if check.Contains(typ.String()) {
			mappingType.SetType(TypeObject)
		} else {
			check.Add(typ.String())

			mappingType.Properties = make(Properties)

			return detectStructProperties(mappingType.Properties, typ, check)
		}
	case reflect.Slice, reflect.Array, reflect.Pointer:
		return detectMappingType(mappingType, typ.Elem(), check)
	case reflect.Map:
		mappingType.SetType(TypeObject)
	case reflect.String:
		mappingType.SetType(TypeText)
	case reflect.Bool:
		mappingType.SetType(TypeBoolean)
	case reflect.Int8:
		mappingType.SetType(TypeByte)
	case reflect.Int16:
		mappingType.SetType(TypeShort)
	case reflect.Int32:
		mappingType.SetType(TypeInteger)
	case reflect.Int, reflect.Int64:
		mappingType.SetType(TypeLong)
	case reflect.Float32:
		mappingType.SetType(TypeFloat)
	case reflect.Float64:
		mappingType.SetType(TypeDouble)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		mappingType.SetType(TypeUnsignedLong)
	default:
		mappingType.SetType(TypeKeyword)
	}

	return nil
}

func detectStructFieldMapping(metadata *MappingType, field reflect.StructField, check collections.Set[string]) (string, error) {
	name, ok := reflection.GetJSONName(field)
	if !ok {
		return "", nil
	}

	if err := metadata.LoadFromTag(field.Tag.Get("es")); err != nil {
		return "", fmt.Errorf("%w: field %s", err, field.Name)
	}

	if metadata.Type == nil {
		if err := detectMappingType(metadata, field.Type, check); err != nil {
			return "", fmt.Errorf("can't detect field '%s' es type: %w", field.Name, err)
		}
	}

	return name, nil
}

func detectStructProperties(properties Properties, typ reflect.Type, check collections.Set[string]) error {
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("%w: can't make mapping for non struct object %v", ErrInvalidMappingObject, typ)
	}

	for i := 0; i < typ.NumField(); i++ {
		metadata := MappingType{}

		property, err := detectStructFieldMapping(&metadata, typ.Field(i), check)
		if err != nil {
			return err
		}

		properties[property] = metadata
	}

	return nil
}
