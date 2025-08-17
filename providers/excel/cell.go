package excel

import (
	"fmt"
	"strconv"
)

// ColumnName возвращает имя колонки excel с заданным индексом. Является отображением множества индексов на множество
// имен колонок листа excel: для 0:"A", 1:"B" и т.д. Для отрицательного индекса вернет "?".
func ColumnName(index int) string {
	if index < 0 {
		return "?"
	}

	n := index / columnCountBase
	if n == 0 {
		return fmt.Sprintf("%c", columnNameA+index)
	}

	return ColumnName(n-1) + ColumnName(index%columnCountBase)
}

// RowName возвращает имя строки excel с заданным индексом. Функция является отображением множества индексов на
// множество имен строк листа excel: для 0:"1", 1:"2" и т.д. Для отрицательного индекса вернет "?".
func RowName(i int) string {
	if i < 0 {
		return "?"
	}

	return strconv.Itoa(i + 1)
}

// CellName возвращает имя ячейки листа excel по заданным индексам колонки col и строки row.
func CellName(col, row int) string {
	return fmt.Sprintf("%s%s", ColumnName(col), RowName(row))
}
