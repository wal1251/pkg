package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/wal1251/pkg/db"
)

// DatabaseExists проверяет существование базы данных.
func DatabaseExists(conn *sql.DB, name string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT schema_name FROM information_schema.schemata WHERE schema_name = ?)"
	if err := conn.QueryRow(query, name).Scan(&exists); err != nil {
		return false, fmt.Errorf("can't query database: %w", err)
	}

	return exists, nil
}

// DatabaseCreate создает новую базу данных. Перед выполнением необходимо выполнить проверку существования вызовом
// DatabaseExists().
func DatabaseCreate(conn *sql.DB, name string) error {
	query := "CREATE DATABASE " + name
	if _, err := conn.Exec(query); err != nil {
		return fmt.Errorf("can't create database %s: %w", name, err)
	}

	return nil
}

// DatabaseDrop удаляет существующую базу данных. Перед выполнением необходимо выполнить проверку существования вызовом
// DatabaseExists().
func DatabaseDrop(conn *sql.DB, name string) error {
	query := "DROP DATABASE " + name
	if _, err := conn.Exec(query); err != nil {
		return fmt.Errorf("can't drop database %s: %w", name, err)
	}

	return nil
}

// DatabaseClear выполнит очистку БД.
func DatabaseClear(conn *sql.DB) error {
	tables, err := TableList(conn)
	if err != nil {
		return err
	}

	err = TableDrop(conn, tables...)
	if err != nil {
		return err
	}

	views, err := ViewList(conn)
	if err != nil {
		return err
	}

	err = ViewDrop(conn, views...)
	if err != nil {
		return err
	}

	return nil
}

// TableList возвращает существующие в БД таблицы, за исключением указанных в exclude.
func TableList(conn *sql.DB, exclude ...string) ([]string, error) {
	dbName, err := currentDBName(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to determine db name: %w", err)
	}

	query := "SELECT table_name FROM information_schema.tables WHERE table_schema = ?"

	args := []interface{}{dbName}

	if len(exclude) > 0 {
		placeholders := strings.Repeat("?,", len(exclude)-1) + "?"
		query += fmt.Sprintf(" AND table_name NOT IN (%s)", placeholders)
		for _, ex := range exclude {
			args = append(args, ex)
		}
	}

	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("can't query database tables: %w", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("can't read query results: %w", err)
		}
		names = append(names, name)
	}

	return names, nil
}

// TableDrop удаляет указанные таблицы БД.
func TableDrop(conn *sql.DB, names ...string) error {
	if len(names) == 0 {
		return nil
	}

	tables := strings.Join(names, ", ")

	query := "DROP TABLE IF EXISTS " + tables
	if _, err := conn.Exec(query); err != nil {
		return fmt.Errorf("can't drop tables %s: %w", tables, err)
	}

	return nil
}

func ViewList(conn *sql.DB, exclude ...string) ([]string, error) {
	dbName, err := currentDBName(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to determine db name: %w", err)
	}

	query := "SELECT table_name FROM information_schema.views WHERE table_schema = ?"

	args := []interface{}{dbName}

	if len(exclude) > 0 {
		placeholders := strings.Repeat("?,", len(exclude)-1) + "?"
		query += fmt.Sprintf(" AND table_name NOT IN (%s)", placeholders)
		for _, ex := range exclude {
			args = append(args, ex)
		}
	}

	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("can't query database views: %w", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("can't read query results: %w", err)
		}
		names = append(names, name)
	}

	return names, nil
}

// ViewDrop удаляет указанные представления БД.
func ViewDrop(conn *sql.DB, names ...string) error {
	if len(names) == 0 {
		return nil
	}

	views := strings.Join(names, ", ")

	query := fmt.Sprintf("DROP VIEW IF EXISTS %s;", views)
	if _, err := conn.Exec(query); err != nil {
		return fmt.Errorf("can't drop views %s: %w", views, err)
	}

	return nil
}

// Truncate удаляет записи таблиц БД, за исключением таблиц, указанных в exclude.
func Truncate(conn *sql.DB, exclude ...string) error {
	exclude = append(exclude, db.SchemaMigrationsTableName)
	names, err := TableList(conn, exclude...)
	if err != nil {
		return err
	}

	// Отключаем проверку внешних ключей
	if _, err = conn.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		return fmt.Errorf("can't disable foreign key checks: %w", err)
	}

	for _, name := range names {
		query := fmt.Sprintf("TRUNCATE TABLE `%s`;", name)
		if _, err = conn.Exec(query); err != nil {
			if _, errSet := conn.Exec("SET FOREIGN_KEY_CHECKS = 1"); errSet != nil {
				return fmt.Errorf("can't truncate table %s: %v, also failed to re-enable foreign key checks: %w", name, err, errSet) // nolint:errorlint
			}

			return fmt.Errorf("can't truncate table %s: %w", name, err)
		}
	}

	// Включаем проверку внешних ключей обратно
	if _, err = conn.Exec("SET FOREIGN_KEY_CHECKS = 1"); err != nil {
		return fmt.Errorf("can't enable foreign key checks: %w", err)
	}

	return nil
}

func currentDBName(conn *sql.DB) (string, error) {
	var currentDB string
	err := conn.QueryRow("SELECT DATABASE()").Scan(&currentDB)
	if err != nil {
		return "", fmt.Errorf("can't get current database name: %w", err)
	}

	return currentDB, nil
}
