package pg_test

import (
	"log"

	"github.com/wal1251/pkg/db"
	"github.com/wal1251/pkg/db/pg"
)

func Example() {
	// Подключение к СУБД.
	conn, err := db.Connect(&db.Config{
		Driver:   db.DriverPostgres,
		Host:     db.CfgDefaultHost,
		Port:     db.CfgDefaultPort,
		SSLMode:  db.SSLModeDisable,
		Database: "postgres",
		User:     "postgres",
		Password: "mysecretpassword",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// БД существует?
	exists, err := pg.DatabaseExists(conn, "sample")
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		// Удалим БД.
		if err = pg.DatabaseDrop(conn, "sample"); err != nil {
			log.Fatal(err)
		}
	}

	// Создадим БД.
	if err = pg.DatabaseCreate(conn, "sample"); err != nil {
		log.Fatal(err)
	}

	// Подключение к новой СУБД sample.
	connSample, err := db.Connect(&db.Config{
		Driver:   db.DriverPostgres,
		Host:     db.CfgDefaultHost,
		Port:     db.CfgDefaultPort,
		SSLMode:  db.SSLModeDisable,
		Database: "sample",
		User:     "postgres",
		Password: "mysecretpassword",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Создадим таблицу foo.
	if _, err = connSample.Exec("CREATE TABLE foo(id varchar primary key, name varchar)"); err != nil {
		log.Fatal(err)
	}

	// Создадим таблицу bar.
	if _, err = connSample.Exec("CREATE TABLE bar(id varchar primary key, name varchar)"); err != nil {
		log.Fatal(err)
	}

	// Запросим таблицы, кроме bar.
	tables, err := pg.TableList(connSample, "bar")
	if err != nil {
		log.Fatal(err)
	}

	// Удалим полученные таблицы.
	if err = pg.TableDrop(connSample, tables...); err != nil {
		log.Fatal(err)
	}

	// Создадим представление bar_view.
	if _, err = connSample.Exec("CREATE VIEW bar_view(id, name) AS SELECT id, name FROM bar"); err != nil {
		log.Fatal(err)
	}

	// Запросим представления, кроме bar.
	views, err := pg.ViewList(connSample, "bar_view")
	if err != nil {
		log.Fatal(err)
	}

	// Удалим полученные представления.
	if err = pg.ViewDrop(connSample, views...); err != nil {
		log.Fatal(err)
	}
}
