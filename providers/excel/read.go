package excel

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"

	"github.com/wal1251/pkg/tools/collections"
)

// SheetRead возвращает считанное с листа sheetName содержимое в collections.Table.
// Если название листа (sheetName) содержит недопустимые символы, вернет ошибку ErrSheetNameInvalid.
// Если лист sheetName не найден, вернет ошибку ErrSheetNotFound. Если имена колонок в таблице excel дублируются, вернет
// ошибку ErrColumnDuplicates. Если задан контроль наличия колонок и нужная колонка не найдена, вернет ошибку
// ErrColumnNotFound.
func SheetRead(file *excelize.File, sheetName string, options ...SheetReadWriteOption) (*collections.Table[string], error) {
	sheetIdx, err := file.GetSheetIndex(sheetName)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrSheetNameInvalid, sheetName)
	}

	if sheetIdx == -1 {
		return nil, fmt.Errorf("%w: %s", ErrSheetNotFound, sheetName)
	}

	readOptions := NewSheetReadWriteOptions(options...)

	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("can't read excel rows from sheet %s: %w", sheetName, err)
	}

	table := collections.NewTable[string]()

	if rows = collections.Skip(rows, readOptions.FirstRowIndex); len(rows) == 0 {
		if readOptions.IsColumnsRequired && len(readOptions.Columns) != 0 {
			return nil, fmt.Errorf("%w: '%s'", ErrColumnNotFound, readOptions.Columns[0])
		}

		return table, nil
	}

	columnNames := collections.NewSet(readOptions.Columns...)
	columnIndexesToSkip := make(collections.Set[int])

	for j, columnName := range collections.Skip(rows[0], readOptions.FirstColumnIndex) {
		columnName = strings.TrimSpace(columnName)
		if columnName == "" || (columnNames.Len() != 0 && !columnNames.Contains(columnName)) {
			columnIndexesToSkip.Add(j)
		} else {
			if collections.NewSet(table.Columns()...).Contains(columnName) {
				return nil, fmt.Errorf("%w: '%s'", ErrColumnDuplicates, columnName)
			}

			table.AddColumns(columnName)
		}
	}

	if readOptions.IsColumnsRequired {
		if columnName, check := collections.NewSet(table.Columns()...).NotContains(readOptions.Columns...); check {
			return nil, fmt.Errorf("%w: '%s'", ErrColumnNotFound, columnName)
		}
	}

	for i := 1; i < len(rows); i++ {
		row := collections.Skip(rows[i], readOptions.FirstColumnIndex)
		table.AddRow(collections.ExceptIndexes(row, columnIndexesToSkip)...)
	}

	return table, nil
}
