package sys

import (
	"fmt"
	"syscall"
	"time"
)

func ExampleShutdownSignal() {
	// Ожидаем сигнал остановки в отдельной горутине.
	go func() {
		<-ShutdownSignal()
		// Здесь можно выполнить необходимую очистку ресурсов
		// и завершить работу приложения.
		fmt.Println("Received shutdown signal")
	}()

	// Имитируем отправку сигнала на завершение приложения.
	go func() {
		time.Sleep(100 * time.Millisecond)
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}()

	time.Sleep(250 * time.Millisecond)

	// Output:
	// Received shutdown signal
}
