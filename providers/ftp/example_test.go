package ftp_test

import (
	"context"
	"fmt"
	"log"

	"github.com/wal1251/pkg/providers/ftp"
	"github.com/wal1251/pkg/tools/files"
)

func ExampleTestDoubleClient() {
	// Пример того, как можно использовать тестовый двойник.

	ctx := context.TODO()

	// Создадим временный файловый том.
	vol := files.NewVolume()
	if err := vol.CreateTemp("fixtures*"); err != nil {
		log.Fatal(err)
	}

	// Обязательно уберем за собой после теста.
	defer func() {
		if err := vol.Remove(); err != nil {
			log.Fatal(err)
		}
	}()

	// Создадим фикстуры в файловом томе.
	if err := files.DirCreateEmptyFiles(vol.Path(),
		"/foo/bar/baz/dummy3.txt",
		"/foo/bar/baz/dummy2.txt",
		"/baz/dummy1.txt",
	); err != nil {
		log.Fatal(err)
	}

	// Инициализируем тестового двойника.
	client := ftp.NewTestDouble(vol.Path())

	// Делаем что-нибудь такое для теста...
	list, err := client.List(ctx, "/foo/bar/baz")
	if err != nil {
		log.Fatal(err)
	}

	// Проверяем...
	for _, entry := range list {
		fmt.Println(entry.Path)
	}

	// Output:
	// /foo/bar/baz/dummy2.txt
	// /foo/bar/baz/dummy3.txt
}
