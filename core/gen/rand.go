package gen

import (
	"math/rand"

	"github.com/wal1251/pkg/core/cfg"
)

// RandomInt возвращает генератор псевдослучайных последовательностей типа int с границей.
func RandomInt() func(int) int {
	if cfg.IsTestRuntime() {
		return DummyRand().Intn
	}

	return rand.Intn
}

// DummyRand возвращает рандомайзер для моков и тестов.
func DummyRand() *rand.Rand {
	return rand.New(rand.NewSource(42)) //nolint:gosec,gomnd
}
