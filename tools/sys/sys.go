package sys

import (
	"os"
	"os/signal"
	"syscall"
)

// ShutdownSignal возвращает канал, который будет получать сигналы
// остановки приложения, такие как os.Interrupt и syscall.SIGTERM.
// Это позволяет приложениям корректно реагировать на сигналы завершения,
// например, для корректного закрытия ресурсов перед завершением работы.
//
// Возвращаемое значение:
// Канал, в который будут отправляться сигналы остановки.
func ShutdownSignal() <-chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	return signals
}
