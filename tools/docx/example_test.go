package docx

import (
	"fmt"
	"os"
	"strings"
)

func ExampleReadDocFromMemory() {
	// Открываем документ.
	file, err := os.Open("testdata/example.docx")
	if err != nil {
		fmt.Println("Error while opening file:", err)
		return
	}
	defer file.Close()

	// Получаем информацию о файле для определения его размера.
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error while getting file info:", err)
		return
	}
	size := fileInfo.Size()

	// Читаем содержимое документа.
	content, err := ReadDocFromMemory(file, size)
	if err != nil {
		fmt.Println("Error while reading doc:", err)
		return
	}

	fmt.Println(strings.Contains(content, "Hello, World!"))

	// Output:
	// true
}
