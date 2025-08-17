package search

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/wal1251/pkg/providers/es/indices"
	"github.com/wal1251/pkg/tools/reflection"
)

const (
	QueryFieldTypeQuery  = "query"
	QueryFieldTypeFilter = "filter"
)

var (
	ErrInvalidObjectType        = errors.New("invalid object type")
	ErrInvalidBoostParam        = errors.New("invalid boost parameter")
	ErrSearchTypeIsNotSpecified = errors.New("search type is not specified")
)

type (
	QueryDescriptor []QueryFieldDescriptor

	QueryFieldDescriptor struct {
		Name  string
		Type  string
		Field string
	}

	DocumentFieldDescriptor struct {
		Name     string
		Type     string
		Field    string
		Boost    *int
		Language *string
	}
)

func (d QueryDescriptor) Visit(query any, accept func(any, QueryFieldDescriptor)) {
	queryValue := reflect.ValueOf(query)
	if queryValue.Type().Kind() == reflect.Pointer {
		if queryValue.IsNil() {
			return
		}
		queryValue = queryValue.Elem()
	}

	for _, field := range d {
		fieldValue := queryValue.FieldByName(field.Field)
		switch fieldValue.Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < fieldValue.Len(); i++ {
				accept(fieldValue.Index(i).Interface(), field)
			}
		case reflect.Pointer:
			if !fieldValue.IsNil() {
				accept(fieldValue.Elem().Interface(), field)
			}
		default:
			accept(fieldValue.Interface(), field)
		}
	}
}

func (d *DocumentFieldDescriptor) LoadFromTag(tagValue string) error {
	if tagValue == "" {
		return nil
	}

	first, pairs := reflection.TagParseKeyValue(tagValue)
	if first != "" {
		d.Type = first
	}

	for key, value := range pairs {
		switch key {
		case "type":
			d.Type = value
		case "name":
			d.Name = value
		case "boost":
			if err := d.parseBoostTag(value); err != nil {
				return err
			}
		case "lang":
			if value != "*" {
				language := value
				d.Language = &language
			}
		}
	}

	return nil
}

func (d *DocumentFieldDescriptor) parseBoostTag(value string) error {
	boostValue, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("%w: failed to parse tag (must be a number): %s", ErrInvalidBoostParam, value)
	}

	if boostValue <= 0 {
		return fmt.Errorf("%w: must be positive number: %s", ErrInvalidBoostParam, value)
	}

	d.Boost = &boostValue

	return nil
}

func (d *QueryFieldDescriptor) LoadFromTag(tagValue string) error {
	if tagValue == "" {
		return nil
	}

	first, pairs := reflection.TagParseKeyValue(tagValue)
	if first != "" {
		d.Type = first
	}

	for key, value := range pairs {
		switch key {
		case "type":
			d.Type = value
		case "name":
			d.Name = value
		}
	}

	if d.Type == "" {
		return ErrSearchTypeIsNotSpecified
	}

	return nil
}

func MakeDocumentFieldDescriptor(name string, mapping indices.MappingType) DocumentFieldDescriptor {
	metadata := DocumentFieldDescriptor{
		Field: name,
		Type:  QueryBoolFilterKey,
	}

	if mapping.Type != nil && *mapping.Type == indices.TypeText {
		metadata.Type = QueryBoolMustKey
	}

	return metadata
}

func MakeQueryFieldDescriptor(name string) QueryFieldDescriptor {
	return QueryFieldDescriptor{
		Field: name,
	}
}

func DocumentFieldsDescription(typ reflect.Type) ([]DocumentFieldDescriptor, error) {
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: can't attain metadata for non struct object %v", ErrInvalidObjectType, typ)
	}

	metadata := make([]DocumentFieldDescriptor, 0)

	for fieldIndex := 0; fieldIndex < typ.NumField(); fieldIndex++ {
		field := typ.Field(fieldIndex)

		mapping, name, err := indices.FieldMapping(field)
		if err != nil {
			return nil, err
		}

		fieldMetadata := MakeDocumentFieldDescriptor(name, mapping)
		if err = fieldMetadata.LoadFromTag(field.Tag.Get("search")); err != nil {
			return nil, fmt.Errorf("%w: field %s", err, name)
		}

		metadata = append(metadata, fieldMetadata)
	}

	return metadata, nil
}

func QueryDescribe(query any) (QueryDescriptor, error) {
	typ := reflect.TypeOf(query)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: can't attain metadata for non struct object %v", ErrInvalidObjectType, typ)
	}

	metadata := make(QueryDescriptor, 0)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldMetadata := MakeQueryFieldDescriptor(field.Name)

		tag := field.Tag.Get("search")
		if tag == "" {
			continue
		}

		if err := fieldMetadata.LoadFromTag(tag); err != nil {
			return nil, fmt.Errorf("%w: field %s", err, field.Name)
		}

		metadata = append(metadata, fieldMetadata)
	}

	return metadata, nil
}
