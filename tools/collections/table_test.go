package collections_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/tools/collections"
)

func TestTable(t *testing.T) {
	table := collections.NewTable[string]("column1", "column2")

	table.AddColumns("column3", "column4")

	require.Equal(t, []string{"column1", "column2", "column3", "column4"}, table.Columns())

	require.Equal(t, 0, table.Size())

	table.AddRow()
	table.AddRow()
	table.AddRow()

	require.Equal(t, 3, table.Size())
}

func TestTableRow_GetValues(t *testing.T) {
	table := collections.NewTable[string]("column1", "column2")

	table.AddRow("1", "2")
	table.AddRow("3", "4")

	require.ElementsMatch(t, []string{"1", "2"}, table.Get(0).GetValues())
	require.ElementsMatch(t, []string{"3", "4"}, table.Get(1).GetValues())
}

func TestTableRow_SetValues(t *testing.T) {
	table := collections.NewTable[string]("column1", "column2")

	table.AddRow()
	require.ElementsMatch(t, []string{"", ""}, table.Get(0).GetValues())

	table.Get(0).SetValues("1", "2")
	require.ElementsMatch(t, []string{"1", "2"}, table.Get(0).GetValues())
}

func TestTableRow_Column(t *testing.T) {
	table := collections.NewTable[string]("column1", "column2")

	table.AddRow("1", "2")
	row := table.Get(0)

	value, ok := row.Column("column1")
	require.True(t, ok)
	require.Equal(t, "1", value)

	value, ok = row.Column("column3")
	require.False(t, ok)
	require.Equal(t, "", value)

	table.AddRow("1")
	row = table.Get(1)

	value, ok = row.Column("column2")
	require.False(t, ok)
	require.Equal(t, "", value)
}

func TestTableRow_ColumnMustPanic(t *testing.T) {
	table := collections.NewTable[string]("column1", "column2")

	table.AddRow("1", "2")
	row := table.Get(0)

	require.Panics(t, func() {
		_ = row.ColumnMust("column3")
	})
}

func TestTableRow_Get(t *testing.T) {
	table := collections.NewTable[string]("column1", "column2")

	table.AddRow("1", "2")
	row := table.Get(0)

	require.Equal(t, row, collections.NewTableRow[string](table, 0))
}

func TestTableRow_GetPanic(t *testing.T) {
	table := collections.NewTable[string]("column1", "column2")

	require.Panics(t, func() {
		_ = table.Get(100)
	})
}

func TestTableEquals(t *testing.T) {
	tests := []struct {
		name  string
		left  func() *collections.Table[int]
		right func() *collections.Table[int]
		want  bool
	}{
		{
			name: "Таблицы равны",
			left: func() *collections.Table[int] {
				tab := collections.NewTable[int]("col1", "col2")
				tab.AddRow(1, 2)
				tab.AddRow(3, 4)
				return tab
			},
			right: func() *collections.Table[int] {
				tab := collections.NewTable[int]("col1", "col2")
				tab.AddRow(1, 2)
				tab.AddRow(3, 4)
				return tab
			},
			want: true,
		},
		{
			name: "Таблицы не равны",
			left: func() *collections.Table[int] {
				tab := collections.NewTable[int]("col1", "col2")
				tab.AddRow(1, 2)
				tab.AddRow(3, 4)
				return tab
			},
			right: func() *collections.Table[int] {
				tab := collections.NewTable[int]("col1", "col2")
				tab.AddRow(0, 1)
				tab.AddRow(3, 4)
				return tab
			},
			want: false,
		},
		{
			name: "Таблицы с разными наименованиями колонок не равны",
			left: func() *collections.Table[int] {
				tab := collections.NewTable[int]("col1", "col2")
				tab.AddRow(1, 2)
				tab.AddRow(3, 4)
				return tab
			},
			right: func() *collections.Table[int] {
				tab := collections.NewTable[int]("col3", "col4")
				tab.AddRow(1, 2)
				tab.AddRow(3, 4)
				return tab
			},
			want: false,
		},
		{
			name: "Пустые таблицы (без колонок)",
			left: func() *collections.Table[int] {
				return collections.NewTable[int]()
			},
			right: func() *collections.Table[int] {
				return collections.NewTable[int]()
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, collections.TableEqual(tt.left(), tt.right()))
		})
	}
}
