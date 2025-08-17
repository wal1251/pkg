package otp

import (
	"github.com/wal1251/pkg/core/memorystore"
	"github.com/wal1251/pkg/providers/otp/generator"
)

// NewManagerTestDouble создает новый экземпляр Manager для тестирования.
// Использует DoubleGenerator для генерации паролей.
func NewManagerTestDouble(sender sender, memoryStore memorystore.MemoryStore, config *Config) *Manager {
	generatorDouble := generator.NewTestDouble(config.Length)

	return &Manager{
		generator:   generatorDouble,
		sender:      sender,
		memoryStore: memoryStore,
		config:      config,
	}
}
