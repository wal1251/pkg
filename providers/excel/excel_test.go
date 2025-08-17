package excel_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/providers/excel"
	"github.com/wal1251/pkg/tools/collections"
)

func TestOpen(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		options []excel.SheetReadWriteOption
		want    *collections.Table[string]
	}{
		{
			name: "Загрузить файл 1",
			file: "./samples_test/sample_3x3.xlsx",
			want: sampleTable(),
		},
		{
			name: "Загрузить файл 2",
			file: "./samples_test/sample_3x3_offset.xlsx",
			options: []excel.SheetReadWriteOption{
				excel.WithFirstRowIndex(1),
				excel.WithFirstColumnIndex(1),
			},
			want: sampleTable(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xl, err := excel.Open(tt.file)
			require.NoError(t, err)
			defer xl.Close(context.TODO())

			sheet, err := xl.Sheet(excel.DefaultSheetName)
			require.NoError(t, err)

			table, err := sheet.Fetch(tt.options...)
			require.NoError(t, err)

			assert.True(t, collections.TableEqual(tt.want, table), "loaded table is not equal to sample")
		})
	}
}

func TestFromReader(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		options []excel.SheetReadWriteOption
		want    *collections.Table[string]
	}{
		{
			name: "Загрузить файл 1",
			file: "./samples_test/sample_3x3.xlsx",
			want: sampleTable(),
		},
		{
			name: "Загрузить файл 2",
			file: "./samples_test/sample_3x3_offset.xlsx",
			options: []excel.SheetReadWriteOption{
				excel.WithFirstRowIndex(1),
				excel.WithFirstColumnIndex(1),
			},
			want: sampleTable(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.file)
			require.NoError(t, err)

			xl, err := excel.FromReader(file)
			require.NoError(t, err)
			defer xl.Close(context.TODO())

			sheet, err := xl.Sheet(excel.DefaultSheetName)
			require.NoError(t, err)

			table, err := sheet.Fetch(tt.options...)
			require.NoError(t, err)

			assert.True(t, collections.TableEqual(tt.want, table), "loaded table is not equal to sample")
		})
	}
}

func sampleTable() *collections.Table[string] {
	table := collections.NewTable[string]("Column_1", "Column_2", "Column_3")
	table.AddRow("Row_1_1", "Row_2_1", "Row_3_1")
	table.AddRow("Row_1_2", "Row_2_2", "Row_3_2")
	table.AddRow("Row_1_3", "Row_2_3", "Row_3_3")
	return table
}
