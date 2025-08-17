package predicates

import (
	"strings"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"

	"github.com/wal1251/pkg/tools/collections"
)

// TuplesIN возвращает предикат-кортеж для выборки значений запросом. Условие строится по table collections.Table
// таким образом:
// WHERE (col1, col2, ...colN) IN (VALUES(row1_col1, row1_col2, ...row1_colN), (row2_col1, row2_col2, ...row2_colN), ...)
//
// Пример:
//
//	tuples := predicates.NewTuples("type", "name")
//	tuples.AddRow("open", "foo")
//	tuples.AddRow("open", "bar")
//
//	if entRefs, err := client.References.Query().Where(predicates.TuplesIN(tuples)).All(ctx); err != nil {
//		return nil, err
//	}
//
// ВНИМАНИЕ! Не работает с nil значениями, т.к. SQL предикат null = null вернет false, не забывайте об этом!
func TuplesIN(table *collections.Table[any]) func(*sql.Selector) {
	columnsTyped := table.Columns()

	return func(selector *sql.Selector) {
		columns := collections.Map(columnsTyped, func(column string) string {
			if selector.Dialect() == dialect.Postgres {
				i := strings.Index(column, "::")
				if i == -1 {
					return column
				}

				return column[:i]
			}

			return column
		})

		if table.Size() == 0 {
			selector.Where(sql.False())

			return
		}

		selector.Where(sql.P(func(builder *sql.Builder) {
			builder.Wrap(func(b *sql.Builder) { b.IdentComma(columns...) }).
				WriteOp(sql.OpIn).
				Wrap(func(builder *sql.Builder) {
					builder.WriteString("SELECT").Pad().
						IdentComma(columnsTyped...).Pad().
						WriteString("FROM").Pad().
						Wrap(func(builder *sql.Builder) {
							builder.WriteString("VALUES").Pad()
							for i := 0; i < table.Size(); i++ {
								if i != 0 {
									builder.Comma()
								}
								builder.Wrap(func(b *sql.Builder) { b.Args(table.Get(i).GetValues()...) })
							}
						}).Pad().
						WriteString("AS").Pad().
						WriteString("__constants").
						Wrap(func(builder *sql.Builder) { builder.IdentComma(columns...) })
				})
		}))
	}
}

// NewTuples создает новый набор кортежей.
func NewTuples(columns ...string) *collections.Table[any] {
	return collections.NewTable[any](columns...)
}
