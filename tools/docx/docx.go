package docx

import (
	"fmt"
	"io"

	dcx "github.com/nguyenthenguyen/docx"
)

// ReadDocFromMemory читает Microsoft Word документ из объекта io.ReaderAt.
//
// Параметры:
// data - объект io.ReaderAt, представляющий данные документа.
// size - размер данных в байтах.
//
// Возвращаемое значение:
// Возвращает содержимое документа в виде строки (в формате xml) и ошибку, если она возникла во время чтения.
func ReadDocFromMemory(data io.ReaderAt, size int64) (string, error) {
	doc, err := dcx.ReadDocxFromMemory(data, size)
	if err != nil {
		return "", fmt.Errorf("failed to read docx from memory: %w", err)
	}

	defer doc.Close()

	content := doc.Editable().GetContent()

	return content, nil
}
