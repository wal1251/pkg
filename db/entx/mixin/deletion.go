package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	"github.com/wal1251/pkg/db/entx"
)

// Deletion реализует ent.Mixin, поля пометки на удаление: флаг удаляемой записи, время установки флага удаления записи.
type Deletion struct {
	mixin.Schema
}

func (Deletion) Fields() []ent.Field {
	return []ent.Field{
		field.Bool(entx.FieldIsDeleted).
			Default(false),
		field.Time(entx.FieldDeletedTime).
			Optional().Nillable(),
	}
}
