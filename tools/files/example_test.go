package files_test

import (
	"fmt"
	"log"
	"os"

	"github.com/wal1251/pkg/tools/files"
)

func Example() {
	// Создадим новый темп каталог.
	samplesDir, err := os.MkdirTemp("", "Samples*")
	if err != nil {
		log.Fatal(err)
	}

	// Приберемся за собой.
	defer func() {
		if err = os.RemoveAll(samplesDir); err != nil {
			log.Fatal(err)
		}
	}()

	// Создадим файлы, кейс может быть полезен для тестов.
	if err = files.DirCreateEmptyFiles(samplesDir,
		"foo/bar/file1",
		"baz/file2",
		"file3",
	); err != nil {
		log.Fatal(err)
	}

	// Создадим новый темп каталог.
	samplesCopyDir, err := os.MkdirTemp("", "SamplesCopy*")
	if err != nil {
		log.Fatal(err)
	}

	// Проверим.
	if err = files.DirMustExist(samplesCopyDir); err != nil {
		// Не должны сюда попасть.
		log.Fatal(err)
	}

	// Приберемся за собой.
	defer func() {
		if err = os.RemoveAll(samplesDir); err != nil {
			log.Fatal(err)
		}
	}()

	// Сделаем копию каталога Samples
	if err = files.DirTraverse(samplesDir, files.DoCopy(samplesCopyDir)); err != nil {
		log.Fatal(err)
	}

	// Выведем перечень скопированных файлов и папок на экран.
	if err = files.DirTraverse(samplesCopyDir, func(entry files.RelativeEntry, skip func()) error {
		if entry.IsHome() {
			return nil
		}

		fmt.Println(entry.PathRelative)

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	// Output:
	// baz
	// baz/file2
	// file3
	// foo
	// foo/bar
	// foo/bar/file1
}

func ExampleVolume() {
	// Создадим новый том.
	vol := files.NewVolume()

	// Привяжем его к временной папке в ФС.
	if err := vol.CreateTemp("Samples*"); err != nil {
		log.Fatal(err)
	}

	// Создадим несколько файлов...
	if err := files.DirCreateEmptyFiles(vol.Path(),
		"file1.txt",
		"foo/bar/file2.txt",
		"baz/file3.txt",
	); err != nil {
		log.Fatal(err)
	}

	// Запросим перечень директорий внутри тома.
	dirs, err := vol.ListDirs()
	if err != nil {
		log.Fatal(err)
	}

	for _, dir := range dirs {
		fmt.Println("FOLDER:", dir)
	}

	// Запросим перечень файлов внутри директории тома.
	list, err := vol.ListFiles("foo")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range list {
		fmt.Println("FILE:", file)
	}

	// Удалим том из ФС.
	if err := vol.Remove(); err != nil {
		log.Fatal(err)
	}

	// Output:
	// FOLDER: .
	// FOLDER: baz
	// FOLDER: foo
	// FOLDER: foo/bar
	// FILE: bar/file2.txt
}
