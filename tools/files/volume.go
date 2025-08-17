package files

import (
	"fmt"
	"os"
	"path"
)

type (
	// Volume файловый том позволяет выполнять различные манипуляции с файловой директорией. Может использоваться в
	// тестах, как источник фикстур.
	Volume struct {
		path string
	}
)

// WithPath установит указанный путь в качестве указателя на директорию ФС.
func (v *Volume) WithPath(newPath string) *Volume {
	v.path = newPath

	return v
}

// Path вернет путь, установленный в качестве указателя на директорию ФС.
func (v *Volume) Path() string {
	return v.path
}

// Load копирует содержимое указанной папки в файловый том.
func (v *Volume) Load(src string) error {
	return DirTraverse(src, DoCopy(v.path))
}

// Create создаст файловый том в ФС, если он не существует.
func (v *Volume) Create() error {
	if err := os.MkdirAll(v.path, DefaultDirPermissions); err != nil {
		return fmt.Errorf("can't create dir %s: %w", v.path, err)
	}

	return nil
}

// CreateTemp создаст новый временный файловый том в ФС.
func (v *Volume) CreateTemp(pattern string) error {
	tmpPath, err := os.MkdirTemp("", pattern)
	if err != nil {
		return fmt.Errorf("can't make tmp dir: %w", err)
	}

	v.path = tmpPath

	return nil
}

// Remove удалит файловый том из ФС.
func (v *Volume) Remove() error {
	if err := os.RemoveAll(v.path); err != nil {
		return fmt.Errorf("can't remove dir %s: %w", v.path, err)
	}

	return nil
}

// ListDirs вернет директории тома.
func (v *Volume) ListDirs() ([]string, error) {
	var list []string
	if err := DirTraverse(v.path, func(entry RelativeEntry, _ func()) error {
		if entry.IsDir() {
			list = append(list, entry.PathRelative)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return list, nil
}

// ListFiles вернет файлы внутри указанной директории тома.
func (v *Volume) ListFiles(basePath string) ([]string, error) {
	var list []string
	if err := DirTraverse(path.Join(v.path, basePath), func(entry RelativeEntry, _ func()) error {
		if !entry.IsDir() {
			list = append(list, entry.PathRelative)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return list, nil
}

// NewVolume создаст новый файловый том.
func NewVolume() *Volume {
	return new(Volume)
}
