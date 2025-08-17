package size

import (
	"bytes"
	"fmt"
	"io"
)

func ExampleCountingWriter() {
	// Создаем буфер в памяти для записи данных.
	buf := bytes.NewBuffer(make([]byte, 0))
	// Создаем экземпляр CountingWriter для отслеживания количества записанных байтов.
	cnt := &CountingWriter{}

	// Используем MultiWriter для записи одновременно в buf и cnt.
	// Это позволяет нам сохранять данные в буфере и отслеживать их объем.
	w := io.MultiWriter(buf, cnt)

	_, err := w.Write([]byte("hello world"))
	if err != nil {
		fmt.Printf("Error during write: %v\n", err)
	}

	fmt.Printf("Total written bytes: %d\n", cnt.Written)

	// Output:
	// Total written bytes: 11
}

func ExampleSize_Count() {
	// Создаем размер в 2 ГБ
	sizeInGB := Make(2, GB)

	// Конвертируем 2 ГБ в Мегабайты
	mbIn2GB := sizeInGB.Count(MB)
	fmt.Printf("2 GB in MB: %d MB", mbIn2GB)
	// Конвертируем 2 ГБ в Килобайты
	kbIn2GB := sizeInGB.Count(KB)
	fmt.Printf("\n2 GB in KB: %d KB", kbIn2GB)

	// Output:
	// 2 GB in MB: 2048 MB
	// 2 GB in KB: 2097152 KB
}

func ExampleSize_String() {
	// Создаем размер в 1024 КБ
	sizeInKB := Make(1024, KB)
	fmt.Println("Size is:", sizeInKB.String())

	// Output:
	// Size is: 1 MB
}

func ExampleBytes() {
	// Конвертируем 1 ГБ в байты
	bytesIn1GB := Bytes(1, GB)
	fmt.Printf("1 GB in bytes: %d bytes", bytesIn1GB)

	// Output:
	// 1 GB in bytes: 1073741824 bytes
}

func ExampleMake() {
	// Создаем размер в 2 ГБ
	sizeInGB := Make(2, GB)
	fmt.Println("Size is:", sizeInGB.String())

	// Output:
	// Size is: 2 GB
}
