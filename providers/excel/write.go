package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"

	"github.com/wal1251/pkg/tools/collections"
)

// SheetWrite записывает содержимое want в лист sheetName. Первая строка считается заголовком таблицы.
func SheetWrite(file *excelize.File, sheetName string, table *collections.Table[string], options ...SheetReadWriteOption) error {
	writeOptions := NewSheetReadWriteOptions(options...)

	for i, columnName := range table.Columns() {
		cell := CellName(writeOptions.FirstColumnIndex+i, writeOptions.FirstRowIndex)
		if err := file.SetCellValue(sheetName, cell, columnName); err != nil {
			return fmt.Errorf("can't set cell %s value in sheet %s: %w", cell, sheetName, err)
		}
	}

	for i := 0; i < table.Size(); i++ {
		r := table.Get(i)
		for j, v := range r.GetValues() {
			cell := CellName(writeOptions.FirstColumnIndex+j, writeOptions.FirstRowIndex+i+1)
			if err := file.SetCellValue(sheetName, cell, v); err != nil {
				return fmt.Errorf("can't set cell %s value in sheet %s: %w", cell, sheetName, err)
			}
		}
	}

	return nil
}
