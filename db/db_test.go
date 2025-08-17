package db_test

import (
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/db"
)

func TestConnect(t *testing.T) {
	conn, err := db.Connect(db.NewCfgSQLiteMem("db"))
	require.NoError(t, err)
	defer conn.Close()

	assert.NotNil(t, conn)
	assert.NoError(t, conn.Ping())
}

func TestLoadFixtures(t *testing.T) {
	// Script1: Создать запись в таблице sample.
	script1, err := os.CreateTemp("", "script1_*.sql")
	require.NoError(t, err)
	defer os.Remove(script1.Name())

	err = os.WriteFile(
		script1.Name(),
		[]byte(`INSERT INTO sample(id,name) VALUES('eecf0477-61e3-4d23-93ad-8260b1fbdb16', 'foo') RETURNING id`),
		os.ModeTemporary,
	)
	require.NoError(t, err)

	// Script2: Создать 2 записи в таблице sample.
	script2, err := os.CreateTemp("", "script2_*.sql")
	require.NoError(t, err)
	defer os.Remove(script2.Name())

	err = os.WriteFile(
		script2.Name(),
		[]byte(`INSERT INTO sample(id,name) VALUES
				('c01e2c6a-043e-4e01-af01-8ea1712626d2', 'bar'),
				('c80e2f92-42ae-4bef-8420-3919f542495f', 'baz')
				RETURNING id;`),
		os.ModeTemporary,
	)
	require.NoError(t, err)

	// Подключение к БД.
	conn, err := db.Connect(db.NewCfgSQLiteMem("db"))
	require.NoError(t, err)
	defer conn.Close()

	// Создадим таблицу для загрузки фикстур.
	_, err = conn.Exec(`CREATE TABLE sample (id TEXT PRIMARY KEY, name TEXT)`)
	require.NoError(t, err)

	// Загружаем фикстуры.
	ids, err := db.LoadFixtures(conn, script1.Name(), script2.Name())
	require.NoError(t, err)

	// Вернули id всех записей?
	assert.ElementsMatch(t, ids, []uuid.UUID{
		uuid.MustParse("eecf0477-61e3-4d23-93ad-8260b1fbdb16"),
		uuid.MustParse("c01e2c6a-043e-4e01-af01-8ea1712626d2"),
		uuid.MustParse("c80e2f92-42ae-4bef-8420-3919f542495f"),
	})

	// Выберем все записи таблицы sample.
	rows, err := conn.Query(`SELECT id FROM sample`)
	require.NoError(t, err)

	dbResult := make([]string, 0, len(ids))

	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		require.NoError(t, err)

		dbResult = append(dbResult, id)
	}

	// В базе id всех записей?
	assert.ElementsMatch(t, dbResult, []string{
		"eecf0477-61e3-4d23-93ad-8260b1fbdb16",
		"c01e2c6a-043e-4e01-af01-8ea1712626d2",
		"c80e2f92-42ae-4bef-8420-3919f542495f",
	})
}
