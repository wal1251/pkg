package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/wal1251/pkg/db"
)

const (
	TypeUUID = "uuid" // Имя типа uuid.
)

// FieldTyped приводит поле к нужному типу.
//
// Например, для сравнения на равенство:
//
//	tuples := predicates.NewTuples(pg.FieldTyped("type", "int"), "name")) // type::int
//	// ...
//
// .
func FieldTyped(field string, typ string) string {
	return fmt.Sprintf("\"%s\"::%s", field, typ)
}

// FieldUUID приводит поле к типу uuid, только для postgres.
//
// Например, для сравнения на равенство:
//
//	tuples := predicates.NewTuples(pg.FieldUUID("type"), "name")) // type::uuid
//	// ...
//
// .
func FieldUUID(field string) string {
	return FieldTyped(field, TypeUUID)
}

// DatabaseExists проверяет существование базы данных.
func DatabaseExists(conn *sql.DB, name string) (bool, error) {
	row := conn.QueryRow(`SELECT true FROM pg_database WHERE datname = $1`, name)
	if err := row.Err(); err != nil {
		return false, fmt.Errorf("can't query database: %w", err)
	}

	var exists bool
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("failed to red query row result: %w", err)
	}

	return exists, nil
}

// DatabaseCreate создает новую базу данных. Перед выполнением необходимо выполнить проверку существования вызовом
// DatabaseExists().
// PG DOCS: "CREATE DATABASE cannot be executed inside a transaction block."
// https://www.postgresql.org/docs/current/sql-createdatabase.html#id-1.9.3.61.7
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
	query := `SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public' AND tablename <> ALL($1)`
	rows, err := conn.Query(query, pq.Array(exclude))
	if err != nil {
		return nil, fmt.Errorf("can't query database tables: %w", err)
	}

	names := make([]string, 0)

	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
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

	query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tables)
	if _, err := conn.Exec(query); err != nil {
		return fmt.Errorf("can't drop tables %s: %w", tables, err)
	}

	return nil
}

// ViewList возвращает существующие в БД представления (views), за исключением указанных в exclude.
func ViewList(conn *sql.DB, exclude ...string) ([]string, error) {
	query := `SELECT table_name FROM INFORMATION_SCHEMA.views WHERE table_schema = 'public' AND table_name <> ALL($1)`
	rows, err := conn.Query(query, pq.Array(exclude))
	if err != nil {
		return nil, fmt.Errorf("can't query database views: %w", err)
	}

	names := make([]string, 0)

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

	query := fmt.Sprintf("DROP VIEW IF EXISTS %s CASCADE;", views)
	if _, err := conn.Exec(query); err != nil {
		return fmt.Errorf("can't drop views %s: %w", views, err)
	}

	return nil
}

// Truncate удаляет записи таблиц БД, за исключением таблиц, указанных в exclude.
func Truncate(conn *sql.DB, exclude ...string) error {
	names, err := TableList(conn, append(exclude, db.SchemaMigrationsTableName)...)
	if err != nil {
		return err
	}

	if len(names) == 0 {
		return nil
	}

	tables := strings.Join(names, ", ")

	query := fmt.Sprintf("TRUNCATE %s CASCADE;", tables)
	if _, err = conn.Exec(query); err != nil {
		return fmt.Errorf("can't truncate tables: %s: %w", query, err)
	}

	return nil
}
