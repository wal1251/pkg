package sentry

import (
	"time"

	"github.com/wal1251/pkg/core/cfg"
)

const (
	CfgDefaultToggle  = true
	CfgDefaultDSN     = ""
	CfgDefaultDebug   = true
	CfgDefaultTimeout = time.Second

	CfgKeySentryToggle  cfg.Key = "SENTRY_TOGGLE"  // Флаг включения/отключения sentry при инициализации проекта
	CfgKeySentryDSN     cfg.Key = "SENTRY_DSN"     // Data Source Name. Адрес отправки событий
	CfgKeySentryDebug   cfg.Key = "SENTRY_DEBUG"   // Флаг включения/отключения режима отладки, вывод дополнительной информации
	CfgKeySentryTimeout cfg.Key = "SENTRY_TIMEOUT" // Время для доставки паники

	tracesSampleRate = 0.2
)

type (
	Config struct {
		Toggle      bool
		DSN         string
		Debug       bool
		Environment string
		Timeout     time.Duration
	}
)
