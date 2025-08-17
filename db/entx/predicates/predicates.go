// Package predicates расширяет возможности стандартной библиотеки ent построения предикатов.
//
// Содержит библиотеку дополнительных sql билдеров ent.
package predicates

import (
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
)

// Optional применяет предикат predicate, только если flag равен true. Удобно задавать предикаты, зависящие от
// вычисляемых значений. Например:
//
//	query.Where(predicates.Optional(folder.ParentID(parentID), parentID != uuid.Nil)
//
// .
func Optional(predicate func(s *sql.Selector), flag bool) func(s *sql.Selector) {
	if flag {
		return predicate
	}

	return func(*sql.Selector) {}
}

// RegexpFold возвращает предикат условия равенства значения колонки column паттерну regexp, без учета регистра
// строки. Не проверяет правильность регулярного выражения.
// Например:
//
//	query.Where(predicates.PgRegexpFold(folder.FieldBoxName, `B\d+`))
//
// ВНИМАНИЕ! Поддерживается только диалект postgres.
func RegexpFold(column string, regexp string) func(s *sql.Selector) {
	return func(s *sql.Selector) {
		s.Where(sql.P(func(b *sql.Builder) {
			if b.Dialect() == dialect.Postgres {
				b.Ident(s.C(column)).Pad().
					WriteString("~*").Pad().
					Arg(regexp)
			}
		}))
	}
}
