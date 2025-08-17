package predicates_test

import (
	"testing"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/db/entx/predicates"
	"github.com/wal1251/pkg/db/pg"
)

func TestTuplesIN(t *testing.T) {
	tuples := predicates.NewTuples("id", "name")
	tuples.AddRow("801a497f-999b-4e9b-bfac-f801d268d0fd", "foo")
	tuples.AddRow("90f65bc1-4284-4f0d-9bf5-d75666dcb6da", "bar")

	builder := sql.Select("id", "name")
	builder.From(sql.Table("sample"))

	condition := predicates.TuplesIN(tuples)
	condition(builder)

	query, args := builder.Query()

	assert.Equal(t, "SELECT `id`, `name` FROM `sample` WHERE (`id`, `name`) IN (SELECT `id`, `name` FROM (VALUES (?, ?), (?, ?)) AS __constants(`id`, `name`))", query)
	assert.Equal(t, 4, len(args))
}

func TestTuplesIN_cast(t *testing.T) {
	tuples := predicates.NewTuples(pg.FieldUUID("id"), "name")
	tuples.AddRow("801a497f-999b-4e9b-bfac-f801d268d0fd", "foo")
	tuples.AddRow("90f65bc1-4284-4f0d-9bf5-d75666dcb6da", "bar")

	builder := sql.Select("id", "name")
	builder.SetDialect(dialect.Postgres)
	builder.From(sql.Table("sample"))

	condition := predicates.TuplesIN(tuples)
	condition(builder)

	query, args := builder.Query()

	assert.Equal(t, `SELECT "id", "name" FROM "sample" WHERE ("id", "name") IN (SELECT "id"::uuid, "name" FROM (VALUES ($1, $2), ($3, $4)) AS __constants("id", "name"))`, query)
	assert.Equal(t, 4, len(args))
}

func TestTuplesIN_empty(t *testing.T) {
	tuples := predicates.NewTuples("id", "name")

	builder := sql.Select("id", "name")
	builder.From(sql.Table("sample"))

	condition := predicates.TuplesIN(tuples)
	condition(builder)

	query, args := builder.Query()

	assert.Equal(t, "SELECT `id`, `name` FROM `sample` WHERE FALSE", query)
	assert.Equal(t, 0, len(args))
}
