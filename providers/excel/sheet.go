package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"

	"github.com/wal1251/pkg/tools/collections"
)

type (
	// Sheet объект, представляющий лист книги excel.
	Sheet struct {
		file  *excelize.File
		index int
		name  string
	}

	// SheetReadWriteOptions опции загрузки\записи excel листа.
	SheetReadWriteOptions struct {
		FirstRowIndex     int      // Индекс первой строки с которой начинается чтение\запись.
		FirstColumnIndex  int      // Индекс первой колонки с которой начинается чтение\запись.
		Columns           []string // Имена колонок, если задан Columns, остальные колонки будут проигнорированы.
		IsColumnsRequired bool     // Проверить, присутствуют ли в таблице колонки SheetReadWriteOptions.Columns.
	}

	SheetReadWriteOption func(*SheetReadWriteOptions)
)

// Exists вернет true, если лист в книге excel существует.
func (s Sheet) Exists() bool {
	return s.index >= 0
}

// Index вернет индекс листа в книге excel.
func (s Sheet) Index() int {
	return s.index
}

// WriteFrom записывает содержимое таблицы в лист excel. Содержимое листа перед записью не очищается.
func (s Sheet) WriteFrom(table *collections.Table[string], options ...SheetReadWriteOption) error {
	return SheetWrite(s.file, s.name, table, options...)
}

// Fetch вернет collections.Table c содержимым листа.
func (s Sheet) Fetch(options ...SheetReadWriteOption) (*collections.Table[string], error) {
	return SheetRead(s.file, s.name, options...)
}

// MakeSheet возвращает лист excel.
// Если название листа (name) содержит недопустимые символы, вернет ошибку ErrSheetNameInvalid.
func MakeSheet(file *excelize.File, name string) (*Sheet, error) {
	sheetIdx, err := file.GetSheetIndex(name)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrSheetNameInvalid, name)
	}

	return &Sheet{
		file:  file,
		name:  name,
		index: sheetIdx,
	}, nil
}

// NewSheetReadWriteOptions возвращает новый SheetReadWriteOptions c установленными опциями.
func NewSheetReadWriteOptions(options ...SheetReadWriteOption) *SheetReadWriteOptions {
	var result SheetReadWriteOptions

	for _, option := range options {
		option(&result)
	}

	return &result
}

// WithFirstRowIndex возвращает функциональную опцию, устанавливающую первую строку для обработки.
func WithFirstRowIndex(index int) SheetReadWriteOption {
	return func(options *SheetReadWriteOptions) {
		options.FirstRowIndex = index
	}
}

// WithFirstColumnIndex возвращает функциональную опцию, устанавливающую первую колонку для обработки.
func WithFirstColumnIndex(index int) SheetReadWriteOption {
	return func(options *SheetReadWriteOptions) {
		options.FirstColumnIndex = index
	}
}

// WithColumns возвращает функциональную опцию, устанавливающую состав колонок для обработки.
func WithColumns(columns ...string) SheetReadWriteOption {
	return func(options *SheetReadWriteOptions) {
		options.Columns = columns
	}
}

// WithColumnsRequired возвращает функциональную опцию, устанавливающую проверку обязательных колонок.
func WithColumnsRequired() SheetReadWriteOption {
	return func(options *SheetReadWriteOptions) {
		options.IsColumnsRequired = true
	}
}
