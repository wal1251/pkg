package ftp

import (
	"fmt"
	"time"

	"github.com/jlaffaye/ftp"
)

const (
	EntryTypeFile   = EntryType(ftp.EntryTypeFile)   // Тип записи - файл.
	EntryTypeFolder = EntryType(ftp.EntryTypeFolder) // Тип записи - папка.
	EntryTypeLink   = EntryType(ftp.EntryTypeLink)   // Тип записи - ссылка.
)

type (
	// Entry запись сервера FTP.
	Entry struct {
		Name   string    // Имя.
		Path   string    // Путь к записи.
		Target string    // Запись, на которую указывает ссылка, актуально для типа EntryTypeLink.
		Type   EntryType // Тип записи.
		Size   uint64    // Размер содержимого.
		Time   time.Time // Время создания записи.
	}

	// EntryType тип записи сервера FTP.
	EntryType int
)

// MakeEntry возвращает новую запись Entry.
func MakeEntry(path string, entry *ftp.Entry) *Entry {
	return &Entry{
		Name:   entry.Name,
		Path:   path,
		Target: entry.Target,
		Type:   EntryType(entry.Type),
		Size:   entry.Size,
		Time:   entry.Time,
	}
}

func (t EntryType) String() string {
	return ftp.EntryType(t).String()
}

// IsFolder вернет true, если запись является папкой.
func (e Entry) IsFolder() bool {
	return e.Type == EntryTypeFolder
}

// IsFile вернет true, если запись является файлом.
func (e Entry) IsFile() bool {
	return e.Type == EntryTypeFile
}

func (e Entry) String() string {
	return fmt.Sprintf("%s %s", e.Type, e.Path)
}
