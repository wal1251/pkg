package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"

	"github.com/wal1251/pkg/core/gen"
	"github.com/wal1251/pkg/db/entx"
)

// ID реализует ent.Mixin, поле с уникальным идентификатором.
type ID struct {
	mixin.Schema
}

func (ID) Fields() []ent.Field {
	return []ent.Field{
		field.UUID(entx.FieldID, uuid.Nil).Default(gen.UUID().Next),
	}
}
