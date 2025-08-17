package indices

import (
	"errors"
	"fmt"

	"github.com/wal1251/pkg/tools/collections"
)

const (
	TypeKeyword      Type = "keyword"
	TypeText         Type = "text"
	TypeObject       Type = "object"
	TypeBoolean      Type = "boolean"
	TypeByte         Type = "byte"
	TypeShort        Type = "short"
	TypeInteger      Type = "integer"
	TypeLong         Type = "long"
	TypeFloat        Type = "float"
	TypeDouble       Type = "double"
	TypeUnsignedLong Type = "unsigned_long"
)

var ErrTypeUnknown = errors.New("type unknown")

type (
	Type string
)

func (t Type) IsValid() bool {
	return collections.NewSet(
		TypeKeyword,
		TypeText,
		TypeObject,
		TypeBoolean,
		TypeByte,
		TypeShort,
		TypeInteger,
		TypeLong,
		TypeFloat,
		TypeDouble,
		TypeUnsignedLong,
	).Contains(t)
}

func parsePropertyType(src string) (*Type, error) {
	typ := Type(src)
	if !typ.IsValid() {
		return nil, fmt.Errorf("%w: %s", ErrTypeUnknown, src)
	}

	return &typ, nil
}
