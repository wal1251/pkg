package files

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

const DefaultDirPermissions = 0o744 // Права по-умолчанию для создаваемых директорий

// DirIsEmpty проверяет, есть ли файлы в указанной директории. Вернет true - если директория пустая.
func DirIsEmpty(dirPath string) (bool, error) {
	file, err := os.Open(dirPath)
	if err != nil {
		return false, fmt.Errorf("can't open dir %s: %w", dirPath, err)
	}

	defer func() {
		_ = file.Close()
	}()

	_, err = file.Readdirnames(1) // Or file.Readdir(1)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return true, nil
		}

		return false, fmt.Errorf("can't read dir %s: %w", dirPath, err)
	}

	return false, nil
}

// DirTraverse рекурсивно обходит содержимое директории dirPath. Для каждой записи ФС вызывает callback. Для пропуска
// сканирования содержимого директории необходимо вызвать skip() в теле callback функции. В качестве пути callback
// функция получает относительный путь файла/папки относительно сканируемого каталога.
func DirTraverse(dirPath string, callback DirVisitor) error {
	err := filepath.WalkDir(dirPath, func(entryPath string, entry fs.DirEntry, _ error) error {
		relativeEntry, err := MakeRelativeEntry(dirPath, entryPath, entry)
		if err != nil {
			return fmt.Errorf("can't calc relative path for %s against %s: %w", entryPath, dirPath, err)
		}

		skip := false
		if err = callback(relativeEntry, func() { skip = true }); err != nil {
			return err
		}

		if skip && entry.IsDir() {
			return filepath.SkipDir
		}

		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrNotExist, dirPath)
		}

		return fmt.Errorf("can't traverse folder %s: %w", dirPath, err)
	}

	return nil
}

// DoCopy возвращает стратегию обхода ФС, которая рекурсивно копирует содержимое директории srcPath в директорию
// destPath. Пути до файлов в целевой директории будут созданы автоматически, если в целевой директории уже присутствует
// целевой файл, вернет ошибку.
func DoCopy(destPath string) DirVisitor {
	return func(entry RelativeEntry, _ func()) error {
		if entry.IsHome() {
			return nil
		}

		destination := entry.Rebase(destPath)

		if entry.IsDir() {
			if err := os.Mkdir(destination, DefaultDirPermissions); err != nil {
				return fmt.Errorf("can't make destination dir %s: %w", destination, err)
			}

			return nil
		}

		destFile, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY, DefaultDirPermissions)
		if err != nil {
			return fmt.Errorf("can't create destination file for copying %s: %w", destination, err)
		}

		srcFile, err := os.Open(entry.Path)
		if err != nil {
			return fmt.Errorf("can't open source file for copying %s: %w", entry.Path, err)
		}

		if _, err = io.Copy(destFile, srcFile); err != nil {
			return fmt.Errorf("failed to copy from %s to %s: %w", entry.Path, destination, err)
		}

		if err = destFile.Close(); err != nil {
			return fmt.Errorf("can't close destination file %s: %w", destination, err)
		}

		if err = srcFile.Close(); err != nil {
			return fmt.Errorf("can't close source file %s: %w", entry.Path, err)
		}

		return nil
	}
}

// DirExists возвращает true, если указанная директория существует в ФС. Вернет ErrIsDir, если указанный путь не
// является директорией.
func DirExists(dirPath string) (bool, error) {
	stat, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, fmt.Errorf("can't read dir entry %s: %w", dirPath, err)
	}

	if !stat.IsDir() {
		return false, fmt.Errorf("%w: %s", ErrIsDir, dirPath)
	}

	return true, nil
}

// DirMustExist проверяет, существует ли в ФС указанная директория. Если не существует, вернет ErrNotExist. Вернет
// ErrIsDir, если указанный путь не является директорией.
func DirMustExist(dirPath string) error {
	ok, err := DirExists(dirPath)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("%w: %s", ErrNotExist, dirPath)
	}

	return nil
}

// DirCreateStructure создает в директории basePath структуру папок dirPaths. Возвращает ошибку, если basePath не
// существует или является файлом.
func DirCreateStructure(basePath string, dirPaths ...string) error {
	if err := DirMustExist(basePath); err != nil {
		return err
	}

	for _, dirPath := range dirPaths {
		newDirPath := path.Join(basePath, dirPath)
		if err := os.MkdirAll(newDirPath, DefaultDirPermissions); err != nil {
			return fmt.Errorf("can't create file path %s: %w", newDirPath, err)
		}
	}

	return nil
}

// DirCreateEmptyFiles создает в директории basePath пустые файлы filePaths. Возвращает ошибку, если basePath не
// существует или является файлом.
// Функция полезна для воссоздания фикстур в ФС для тестов.
func DirCreateEmptyFiles(basePath string, filePaths ...string) error {
	if err := DirMustExist(basePath); err != nil {
		return err
	}

	for _, filePath := range filePaths {
		fullPath := path.Join(basePath, filePath)
		fileDirPath, _ := filepath.Split(fullPath)

		if err := os.MkdirAll(fileDirPath, DefaultDirPermissions); err != nil {
			return fmt.Errorf("can't create file path %s: %w", fileDirPath, err)
		}

		file, err := os.Create(fullPath)
		if err != nil {
			return fmt.Errorf("can't create file %s: %w", fullPath, err)
		}

		if err = file.Close(); err != nil {
			return fmt.Errorf("can't close file %s: %w", fullPath, err)
		}
	}

	return nil
}
