package redis

import (
	"fmt"

	"github.com/alicebob/miniredis/v2"
)

// TestRedisServer представляет собой тестовый сервер Redis.
type TestRedisServer struct {
	double *miniredis.Miniredis
}

// NewTestRedisServer создает и возвращает новый экземпляр TestRedisServer.
func NewTestRedisServer() TestRedisServer {
	return TestRedisServer{
		double: miniredis.NewMiniRedis(),
	}
}

// Run запускает тестовый сервер Redis на заданном адресе и порту, указанных в конфигурации.
func (tc *TestRedisServer) Run(cfg Config) error {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	if err := tc.double.StartAddr(addr); err != nil {
		return fmt.Errorf("can't start mock redis server on %s: %w", addr, err)
	}

	return nil
}

// FlushAll удаляет все данные из тестового сервера Redis.
func (tc *TestRedisServer) FlushAll() {
	tc.double.FlushAll()
}

// Close останавливает и закрывает тестовый сервер Redis, освобождая все ресурсы,
// занятые внутренним инстансом miniredis.
func (tc *TestRedisServer) Close() {
	tc.double.Close()
}
