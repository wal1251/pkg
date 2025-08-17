package logs

import "github.com/wal1251/pkg/core/cfg"

const (
	LevelInfo = "info" // Уровень логирования 'info'.

	CfgKeyLevel  cfg.Key = "LOG_LEVEL"  // Конфиг: string - уровень логирования.
	CfgKeyPretty cfg.Key = "LOG_PRETTY" // Конфиг: bool - форматированный вывод логов.

	CfgDefaultLevel = LevelInfo // Уровень логирования по умолчанию.
)

// Config параметры конфигурации логера.
type Config struct {
	Level  string // Уровень логирования.
	Pretty bool   // Форматированный вывод логов.
}
