package excel_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/providers/excel"
)

func TestColumnName(t *testing.T) {
	tests := []struct {
		name string
		i    int
		want string
	}{
		{
			name: "Колонка A",
			i:    0,
			want: "A",
		},
		{
			name: "Колонка D",
			i:    3,
			want: "D",
		},
		{
			name: "Колонка Z",
			i:    25,
			want: "Z",
		},
		{
			name: "Колонка AA",
			i:    26,
			want: "AA",
		},
		{
			name: "Колонка AD",
			i:    29,
			want: "AD",
		},
		{
			name: "Колонка AMJ",
			i:    1023,
			want: "AMJ",
		},
		{
			name: "Отрицательный индекс",
			i:    -1,
			want: "?",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, excel.ColumnName(tt.i))
		})
	}
}

func TestRowName(t *testing.T) {
	tests := []struct {
		name string
		i    int
		want string
	}{
		{
			name: "Строка 1",
			i:    0,
			want: "1",
		},
		{
			name: "Строка 4",
			i:    3,
			want: "4",
		},
		{
			name: "Строка 1024",
			i:    1023,
			want: "1024",
		},
		{
			name: "Отрицательный индекс",
			i:    -1,
			want: "?",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, excel.RowName(tt.i))
		})
	}
}

func TestCellName(t *testing.T) {
	tests := []struct {
		name string
		col  int
		row  int
		want string
	}{
		{
			name: "Ячейка A1",
			col:  0,
			row:  0,
			want: "A1",
		},
		{
			name: "Ячейка AMJ65536",
			col:  1023,
			row:  65535,
			want: "AMJ65536",
		},
		{
			name: "Отрицательный индекс",
			col:  -1,
			row:  -1,
			want: "??",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, excel.CellName(tt.col, tt.row))
		})
	}
}
