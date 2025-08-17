// Package files предоставляет вспомогательные хелперы для работы с файлами.
package files

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

var (
	ErrIsDir    = errors.New("is directory") // Является директорией.
	ErrNotExist = errors.New("not exist")    // Не существует.
)

type (
	RelativeEntry struct {
		os.DirEntry

		Path         string
		PathBase     string
		PathRelative string
	}

	// DirVisitor реализует стратегию обхода папок/файлов для файловой системы. Для отказа от обхода вглубь, в теле
	// визитёра необходимо вызвать skip(). Если DirVisitor возвращает error, тогда вызывающая функция возвращает тот же
	// самый error.
	DirVisitor func(entry RelativeEntry, skip func()) error
)

func MakeRelativeEntry(basePath, absPath string, entry os.DirEntry) (RelativeEntry, error) {
	relPath, err := filepath.Rel(basePath, absPath)
	if err != nil {
		return RelativeEntry{}, fmt.Errorf("can't calc relative path for %s against %s: %w", absPath, basePath, err)
	}

	return RelativeEntry{
		Path:         absPath,
		PathBase:     basePath,
		PathRelative: relPath,
		DirEntry:     entry,
	}, nil
}

func (e RelativeEntry) IsHome() bool {
	return e.IsDir() && e.PathRelative == "."
}

func (e RelativeEntry) Rebase(basePath string) string {
	return path.Join(basePath, e.PathRelative)
}
