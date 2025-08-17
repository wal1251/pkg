package collections

import (
	"errors"
	"fmt"

	"golang.org/x/exp/slices"
)

var (
	ErrUnknownColumnName = errors.New("unknown column name")
	ErrIndexOutOfRange   = errors.New("index out of range")
)

type (
	// Table таблица строк, значения столбцов которой представлены типом T.
	Table[T any] struct {
		columns map[string]int
		rows    [][]T
	}

	// TableRow строка таблицы Table.
	TableRow[T any] struct {
		table *Table[T]
		index int
	}
)

// Columns получить список столбцов в таблице.
func (t *Table[T]) Columns() []string {
	columns := make([]string, len(t.columns))

	for name, i := range t.columns {
		columns[i] = name
	}

	return columns
}

// AddColumns добавить столбец.
func (t *Table[T]) AddColumns(names ...string) {
	for _, name := range names {
		if _, exists := t.columns[name]; exists {
			continue
		}

		t.columns[name] = len(t.columns)
	}
}

// AddRow добавить строку.
func (t *Table[T]) AddRow(values ...T) *TableRow[T] {
	index := len(t.rows)
	t.rows = append(t.rows, []T{})

	row := &TableRow[T]{
		table: t,
		index: index,
	}

	row.SetValues(values...)

	return row
}

// Get Получить TableRow по индексу.
// Если индекс больше размера таблицы, то будет паника.
func (t *Table[T]) Get(index int) *TableRow[T] {
	return NewTableRow(t, index)
}

// Size получить размер таблицы(количество строк).
func (t *Table[T]) Size() int {
	return len(t.rows)
}

// ColumnMust получить значение в строке TableRow и в столбце с name.
// Если такой столбец отсутствует, то будет паника.
func (r *TableRow[T]) ColumnMust(name string) T {
	if v, ok := r.Column(name); ok {
		return v
	}

	panic(fmt.Errorf("%w: %s", ErrUnknownColumnName, name))
}

// Column получить значение в строке TableRow и в столбце с name.
// Если такой столбец отсутствует, то вернуть дефолтное значение типа T и false.
func (r *TableRow[T]) Column(name string) (T, bool) {
	var zero T

	if columnIndex, ok := r.table.columns[name]; ok {
		rowData := r.table.rows[r.index]

		if columnIndex >= len(rowData) {
			return zero, false
		}

		return rowData[columnIndex], true
	}

	return zero, false
}

// SetValues установить список значений в TableRow.
func (r *TableRow[T]) SetValues(values ...T) {
	r.table.rows[r.index] = values
}

// GetValues получить список значений в TableRow.
func (r *TableRow[T]) GetValues() []T {
	columnsCount := len(r.table.columns)

	values := r.table.rows[r.index]
	if len(values) != columnsCount {
		normalized := make([]T, columnsCount)

		copy(normalized, values)

		values = normalized
	}

	return values
}

// Index получить индекс TableRow.
func (r *TableRow[T]) Index() int {
	return r.index
}

// NewTable возвращает новую Table c колонками типа T.
func NewTable[T any](columns ...string) *Table[T] {
	table := &Table[T]{
		columns: make(map[string]int),
		rows:    make([][]T, 0),
	}

	table.AddColumns(columns...)

	return table
}

// NewTableRow возвращает новую TableRow, привязанную к таблице Table, c колонками типа T.
// Если индекс больше размера таблицы, то будет паника.
func NewTableRow[T any](table *Table[T], index int) *TableRow[T] {
	if index > table.Size()-1 {
		panic(fmt.Errorf("%w: %d", ErrIndexOutOfRange, index))
	}

	return &TableRow[T]{
		table: table,
		index: index,
	}
}

// TableEqual вернет true, если таблицы равны (наименования колонок и содержимое строк). Сложность в худшем случае
// (равенство таблиц) составит O(N*M).
func TableEqual[T comparable](opl, opr *Table[T]) bool {
	if opl.Size() != opr.Size() || !slices.Equal(opl.Columns(), opr.Columns()) {
		return false
	}

	for i := 0; i < opr.Size(); i++ {
		if !slices.Equal(opl.Get(i).GetValues(), opr.Get(i).GetValues()) {
			return false
		}
	}

	return true
}
