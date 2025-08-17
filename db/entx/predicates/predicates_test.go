package predicates_test

import (
	"testing"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/db/entx/predicates"
)

func TestOptional(t *testing.T) {
	falsePredicate := func(s *sql.Selector) { s.Where(sql.False()) }

	tests := []struct {
		name string
		flag bool
		want string
	}{
		{
			name: "Flag is true",
			flag: true,
			want: "SELECT `id`, `name` FROM `sample` WHERE FALSE",
		},
		{
			name: "Flag is false",
			flag: false,
			want: "SELECT `id`, `name` FROM `sample`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := sql.Select("id", "name")
			builder.From(sql.Table("sample"))

			condition := predicates.Optional(falsePredicate, tt.flag)
			condition(builder)

			query, _ := builder.Query()
			assert.Equal(t, tt.want, query)
		})
	}
}

func TestRegexpFold(t *testing.T) {
	tests := []struct {
		name    string
		dialect string
		column  string
		regexp  string
		want    string
	}{
		{
			name:    "Базовый кейс",
			dialect: dialect.Postgres,
			column:  "name",
			regexp:  `\w+`,
			want:    "SELECT \"id\", \"name\" FROM \"sample\" WHERE \"sample\".\"name\" ~* $1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := sql.Select("id", "name")
			builder.From(sql.Table("sample"))
			builder.SetDialect(tt.dialect)

			condition := predicates.RegexpFold(tt.column, tt.regexp)
			condition(builder)

			query, args := builder.Query()
			assert.Equal(t, tt.want, query)
			assert.Equal(t, 1, len(args), "количество ожидаемых аргументов не равно 1")
		})
	}
}
