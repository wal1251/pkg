package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file" // register file migration source.
	"github.com/rs/zerolog"

	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/tools/files"
)

var _ migrate.Logger = (*MigrateLogger)(nil)

const SchemaMigrationsTableName = "schema_migrations"

var ErrDriverIsUnsupported = errors.New("driver is unsupported") // Указанный драйвер не поддерживается.

type (
	// MigrateLogger специализированная реализация логера, используемый migrate. По сути, адаптер, позволяющий migrate
	// вызывать стандартный логер. Реализует migrate.Logger.
	MigrateLogger struct {
		Logger    *zerolog.Logger // Целевой логер.
		IsVerbose bool            // Подробный вывод логов.
	}
)

// Printf см. migrate.Logger.Printf().
func (m MigrateLogger) Printf(format string, v ...any) {
	m.Logger.Debug().
		Str(string(logs.ComponentTag), "db.Migrate").
		Msgf(strings.TrimRight(format, "\n"), v...)
}

// Verbose см. migrate.Logger.Verbose().
func (m MigrateLogger) Verbose() bool {
	return m.IsVerbose
}

// Migrate выполняет миграцию, согласно скриптам источника, указанного в конфиге. Требует указания конкретного драйвера
// миграции. Но можно попытаться определить драйвер автоматически на основе конфига, см. MigrateDefault().
func Migrate(ctx context.Context, cfg ConnectionDescriber, newDriver func(*sql.DB) (database.Driver, error)) error {
	logger := logs.FromContext(ctx)

	if cfg.MigrationSource() == "" {
		logger.Info().Msg("no DB migrations source is specified: migrations skipped")

		return nil
	}

	if _, err := os.Stat(cfg.MigrationSource()); err != nil {
		return fmt.Errorf("can't get DB migrations dir (%s): %w", cfg.MigrationSource(), err)
	}

	isEmptyDir, err := files.DirIsEmpty(cfg.MigrationSource())
	if err != nil {
		return fmt.Errorf("can't lookup DB migrations dir (%s): %w", cfg.MigrationSource(), err)
	}

	if isEmptyDir {
		logger.Info().Msg("DB migrations source is empty: migrations skipped")

		return nil
	}

	logger.Info().Msgf("start applying DB migrations from source: %s", cfg.MigrationSource())

	conn, err := Connect(cfg)
	if err != nil {
		return err
	}

	driver, err := newDriver(conn)
	if err != nil {
		return err
	}

	source := "file://" + cfg.MigrationSource()
	migrateInstance, err := migrate.NewWithDatabaseInstance(source, cfg.DatabaseName(), driver)
	if err != nil {
		return fmt.Errorf("cannot get migrate instance: %w", err)
	}

	migrateInstance.Log = MigrateLogger{
		Logger:    logs.FromContext(ctx),
		IsVerbose: cfg.IsDebug(),
	}

	if err = migrateInstance.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migration error: %w", err)
		}
	}

	return nil
}

// MigratePostgres предоставляет конструктор драйвера миграции для PostgresSQL.
func MigratePostgres(db *sql.DB) (database.Driver, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("can't get migration driver instance: %w", err)
	}

	return driver, nil
}

// MigrateMySQL предоставляет конструктор драйвера миграции для MySQL.
func MigrateMySQL(db *sql.DB) (database.Driver, error) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, fmt.Errorf("can't get migration driver instance: %w", err)
	}

	return driver, nil
}

// MigrateSQLite3 предоставляет конструктор драйвера миграции для SQLite.
func MigrateSQLite3(db *sql.DB) (database.Driver, error) {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return nil, fmt.Errorf("can't get migration driver instance: %w", err)
	}

	return driver, nil
}

// MigrateDefault возвращает подходящий конструктор драйвера на основе конфигурации подключения.
func MigrateDefault(connDescriber ConnectionDescriber) func(*sql.DB) (database.Driver, error) {
	switch connDescriber.DriverName() {
	case DriverPostgres:
		return MigratePostgres
	case DriverMySQL:
		return MigrateMySQL
	case DriverSQLite3:
		return MigrateSQLite3
	default:
		return func(*sql.DB) (database.Driver, error) {
			return nil, fmt.Errorf("%w: %s", ErrDriverIsUnsupported, connDescriber.DriverName())
		}
	}
}
