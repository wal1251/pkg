package mixin

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	"github.com/wal1251/pkg/db/entx"
)

// Time реализует ent.Mixin, поля временных меток: время создания записи, время обновления записи.
type Time struct {
	mixin.Schema
}

func (Time) Fields() []ent.Field {
	return []ent.Field{
		field.Time(entx.FieldCreateTime).
			Default(time.Now).
			Immutable(),
		field.Time(entx.FieldUpdateTime).
			Default(time.Now).UpdateDefault(time.Now).
			Immutable(),
	}
}
