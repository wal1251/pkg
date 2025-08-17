// Package db содержит объекты и функции для работы с реляционными базами данных.
package db

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"github.com/wal1251/pkg/tools/collections"
)

var (
	ErrNotMySQLDriver    = errors.New("not a MySQL driver")
	ErrFailedToAppendPEM = errors.New("failed to append PEM")
)

// Connect возвращает открытое соединение с БД и проверяет его. Если во время проверки произошла ошибка, вернет error.
func Connect(cfg ConnectionDescriber) (*sql.DB, error) {
	if cfg.DriverName() == DriverMySQL {
		err := RegisterMySQLTLSConfig(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to register tls config for MySQL: %w", err)
		}
	}

	db, err := sql.Open(cfg.DriverName(), cfg.DataSourceName())
	if err != nil {
		return nil, fmt.Errorf("failed to open %s connection (%v): %w", cfg.DriverName(), cfg, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("can't check db connection (%v): %w", cfg, err)
	}

	return db, nil
}

// DefaultConnect возвращает открытое соединение с БД и проверяет его. Если во время проверки произошла ошибка, вернет error.
// Данный метод можно использовать для создания или удаления БД, когда инициализируешь тестовое окружение.
func DefaultConnect(cfg ConnectionDescriber) (*sql.DB, error) {
	db, err := sql.Open(cfg.DriverName(), cfg.DefaultDataSourceName())
	if err != nil {
		return nil, fmt.Errorf("failed to open %s connection (%v): %w", cfg.DriverName(), cfg, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("can't check db connection (%v): %w", cfg, err)
	}

	return db, nil
}

// RegisterMySQLTLSConfig регистрирует конфигурацию TLS для MySQL
// для режимов SSL, где необходим сертификат (SSLModeVerifyCA, SSLModeVerifyFull).
// Эта функция должна быть вызвана до использования sql.Open() с DSN,
// который был получен через ConnectionDescriber.DataSourceName().
//
// Вызов необходим, поскольку драйвер MySQL не поддерживает указание пути к SSL сертификату
// непосредственно в DSN, в отличие от драйвера для Postgres.
// Подробности смотрите в issue: https://github.com/go-sql-driver/mysql/issues/926
func RegisterMySQLTLSConfig(cfg ConnectionDescriber) error {
	if cfg.DriverName() != DriverMySQL {
		return fmt.Errorf("%w: current driver - %s", ErrNotMySQLDriver, cfg.DriverName())
	}

	tlsConfigRegisterNotRequire := collections.NewSet(SSLModeDisable, SSLModeRequire)
	if tlsConfigRegisterNotRequire.Contains(cfg.SSLConnectionMode()) {
		return nil
	}

	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(cfg.SSLCertificatePath())
	if err != nil {
		return fmt.Errorf("failed to read ssl cert file: %w", err)
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return ErrFailedToAppendPEM
	}

	tlsConfig := &tls.Config{
		RootCAs:    rootCertPool,
		MinVersion: tls.VersionTLS12,
	}
	if cfg.SSLConnectionMode() == SSLModeVerifyFull {
		tlsConfig.ServerName = cfg.DatabaseHost()
	}

	// Регистрируем конфигурацию TLS используя в качестве ключа конфига имя SSL.
	// Ключ конфига ранее задан, как параметр tls в DSN, см. ConnectionDescriber.DataSourceName().
	err = mysql.RegisterTLSConfig(cfg.SSLConnectionMode(), tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to register tls config: %w", err)
	}

	return nil
}

// LoadFixtures загружает фикстуры из указанных файлов и возвращает id созданных записей.
func LoadFixtures(conn *sql.DB, fileNames ...string) ([]uuid.UUID, error) {
	result := make([]uuid.UUID, 0, len(fileNames))

	for _, file := range fileNames {
		script, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file (%s): %w", file, err)
		}

		rows, err := conn.Query(string(script))
		if err != nil {
			return nil, fmt.Errorf("can't query db: %w", err)
		}

		for rows.Next() {
			var id string
			if err = rows.Scan(&id); err != nil {
				return nil, fmt.Errorf("can't read query results: %w", err)
			}

			guid, err := uuid.Parse(id)
			if err != nil {
				return nil, fmt.Errorf("id value is incorrect: %w", err)
			}

			result = append(result, guid)
		}
	}

	return result, nil
}
