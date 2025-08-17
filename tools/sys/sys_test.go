package sys_test

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/tools/sys"
)

func TestShutdownSignal(t *testing.T) {
	signals := sys.ShutdownSignal()

	// Имитируем отправку сигнала на завершение приложения.
	go func() {
		time.Sleep(100 * time.Millisecond)
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}()

	select {
	case sig := <-signals:
		require.Equalf(t, syscall.SIGTERM, sig, "Expected SIGTERM, got %s", sig)
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Timeout waiting for shutdown signal")
	}
}
