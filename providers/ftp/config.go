package ftp

import (
	"time"

	"github.com/wal1251/pkg/core/cfg"
)

const (
	CfgKeyAddress  cfg.Key = "FTP_ADDRESS"  // Ключ конфигурации: адрес сервера (string).
	CfgKeyUser     cfg.Key = "FTP_USER"     // Ключ конфигурации: пользователь для аутентификации на сервере (string).
	CfgKeyPassword cfg.Key = "FTP_PASSWORD" // Ключ конфигурации: пароль пользователя для аутентификации на сервере (string).
	CfgKeyTimeout  cfg.Key = "FTP_TIMEOUT"  // Ключ конфигурации: таймаут клиента FTP (duration).

	CfgDefaultAddress  = "localhost:21"  // Адрес сервера по-умолчанию.
	CfgDefaultUser     = "anonymous"     // Пользователь по-умолчанию.
	CfgDefaultPassword = "guest"         // Пароль пользователя по-умолчанию.
	CfgDefaultTimeout  = 5 * time.Second // Таймаут клиента по-умолчанию.
)

// Config конфигурация клиента, для подключения к серверу FTP.
type Config struct {
	Address  string        // Адрес сервера.
	User     string        // Пользователь для аутентификации на сервере.
	Password string        // Пароль пользователя для аутентификации на сервере.
	Timeout  time.Duration // Таймаут клиента FTP.
}
