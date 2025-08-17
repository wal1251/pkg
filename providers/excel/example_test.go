package excel_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/wal1251/pkg/providers/excel"
	"github.com/wal1251/pkg/tools/collections"
)

func ExampleExcel() {
	// Подготовим данные для записи в лист excel.
	table := collections.NewTable[string]("C.1", "C.2", "C.3")
	table.AddRow("R.1.1", "R.2.1", "R.3.1")
	table.AddRow("R.1.2", "R.2.2", "R.3.2")
	table.AddRow("R.1.3", "R.2.3", "R.3.3")

	ctx := context.TODO()

	// Создадим новую книгу excel.
	xlWriter := excel.New()
	defer xlWriter.Close(ctx)

	// Создадим новый лист MySheet и запишем в него ранее подготовленную информацию.
	sheet, err := xlWriter.SheetCreate("MySheet")
	if err != nil {
		log.Fatal(err)
	}

	if err := sheet.WriteFrom(table); err != nil {
		log.Fatal(err)
	}

	// Создадим временный файл, в него будем писать сформированную книгу.
	file, err := os.CreateTemp("", "sample_*.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	// Наша ответственность - убрать за собой.
	defer func() {
		if err = os.Remove(file.Name()); err != nil {
			log.Fatal(err)
		}
	}()

	// Запишем книгу excel в файл.
	if _, err = xlWriter.WriteTo(file); err != nil {
		log.Fatal(err)
	}

	// Закроем сейчас, т.к. хотим из файла читать.
	if err = file.Close(); err != nil {
		log.Fatal(err)
	}

	// Прочитаем только что созданную книгу excel из ФС.
	xlReader, err := excel.Open(file.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer xlReader.Close(ctx)

	// Прочитаем лист книги excel.
	sheet, err = xlReader.SheetCreate("MySheet")
	if err != nil {
		log.Fatal(err)
	}

	tableRead, err := sheet.Fetch()
	if err != nil {
		log.Fatal(err)
	}

	// То, что мы считали из ФС совпадает с первоначальными данными, которые мы хотели разместить на листе?
	if collections.TableEqual(table, tableRead) {
		fmt.Println("OK")
	}

	// Output:
	// OK
}
