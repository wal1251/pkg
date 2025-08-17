// Package excel предоставляет объекты и функции для работы с таблицами в формате excel: чтение, создание, модификация.
package excel

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"

	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/tools/collections"
)

const (
	columnNameA     = 'A'           // Имя колонки A (первая колонка).
	columnCountBase = 'Z' - 'A' + 1 // Основание системы исчисления для нумерации колонок excel.

	DefaultSheetName = "Sheet1" // Иля листа excel по-умолчанию.
)

var (
	ErrSheetNotFound    = errors.New("sheet not found")    // Лист excel не найден.
	ErrSheetNameInvalid = errors.New("sheet name invalid") // Невалидное название для листа excel.
	ErrColumnDuplicates = errors.New("column duplicates")  // Колонка в таблице дублируется.
	ErrColumnNotFound   = errors.New("column not found")   // Не найдена обязательная колонка.
)

var _ io.WriterTo = Excel{}

type (
	// Excel предоставляет методы для работы с книгой excel.
	Excel struct {
		file *excelize.File
	}
)

// Reader пишет содержимое книги excel во временный буфер и возвращает io.Reader для считывания буфера. Не подходит для
// больших файлов, т.к. размещает содержимое файла в памяти целиком. В случае больших файлов предпочтительнее
// использовать Excel.WriteTo().
func (e Excel) Reader() (io.Reader, error) {
	buffer, err := e.file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("can't write excel book: %w", err)
	}

	return buffer, nil
}

// WriteTo предоставляет реализацию стандартного интерфейса io.WriterTo.
func (e Excel) WriteTo(w io.Writer) (int64, error) {
	cnt, err := e.file.WriteTo(w)
	if err != nil {
		return 0, fmt.Errorf("can't write excel book: %w", err)
	}

	return cnt, nil
}

// Close высвобождает занятые ресурсы.
func (e Excel) Close(ctx context.Context) {
	logger := logs.FromContext(ctx)
	if err := e.file.Close(); err != nil {
		logger.Err(err).Msg("failed to close excel book")
	}
}

// Sheet получает лист книги excel по имени. Например:
//
//	xl := excel.New()
//	sheet, err := xl.Sheet(excel.DefaultSheetName)
//	if err != nil {
//		// невалидное имя листа ...
//		return
//	}
//
//	if !sheet.Exists() {
//		// такого листа в книге нет ...
//		return
//	}
//
//	// работаем с листом ...
//
// .
func (e Excel) Sheet(name string) (*Sheet, error) {
	return MakeSheet(e.file, name)
}

// Sheets вернет все листы в книге excel.
func (e Excel) Sheets() []*Sheet {
	return collections.Map(e.file.GetSheetList(), func(name string) *Sheet {
		sheet, _ := MakeSheet(e.file, name)

		return sheet
	})
}

// SheetCreate создает и возвращает новый лист excel.
// Если название листа (name) содержит недопустимые символы, вернет ошибку ErrSheetNameInvalid.
func (e Excel) SheetCreate(name string) (*Sheet, error) {
	sheetIdx, err := e.file.NewSheet(name)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrSheetNameInvalid, name)
	}

	return &Sheet{
		file:  e.file,
		name:  name,
		index: sheetIdx,
	}, nil
}

// SheetGetOrCreate возвращает лист excel, если лист с таким именем уже существует, в противном случае создает новый.
// Если название листа (name) содержит недопустимые символы, вернет ошибку ErrSheetNameInvalid.
func (e Excel) SheetGetOrCreate(name string) (*Sheet, error) {
	target, err := MakeSheet(e.file, name)
	if err != nil {
		return nil, err
	}

	if target.Exists() {
		return target, nil
	}

	return e.SheetCreate(name)
}

// SheetRename изменяет имя листа c currentName на newName.
func (e Excel) SheetRename(currentName, newName string) error {
	err := e.file.SetSheetName(currentName, newName)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrSheetNameInvalid, newName)
	}

	return nil
}

// New создает и возвращает новый объект книги excel.
func New() Excel {
	return Excel{file: excelize.NewFile()}
}

// FromReader возвращает объект книги excel со считанными в него содержимым из io.Reader.
func FromReader(reader io.Reader, opts ...excelize.Options) (Excel, error) {
	f, err := excelize.OpenReader(reader, opts...)
	if err != nil {
		return Excel{}, fmt.Errorf("can't read excel: %w", err)
	}

	return Excel{file: f}, nil
}

// Open возвращает объект книги excel со считанными в него содержимым из файла с заданным именем.
func Open(name string, opts ...excelize.Options) (Excel, error) {
	f, err := excelize.OpenFile(name, opts...)
	if err != nil {
		return Excel{}, fmt.Errorf("can't read excel file %s: %w", name, err)
	}

	return Excel{file: f}, nil
}
