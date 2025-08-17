package mysql_test

import (
	"log"

	"github.com/wal1251/pkg/db"
	"github.com/wal1251/pkg/db/mysql"
)

func Example() {
	// Подключение к СУБД.
	conn, err := db.Connect(&db.Config{
		Driver:   db.DriverMySQL,
		Host:     db.CfgDefaultHost,
		Port:     "3306",
		SSLMode:  db.SSLModeDisable,
		Database: "mysql",
		User:     "root",
		Password: "secret",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	dbName := "sample"
	// БД существует?
	exists, err := mysql.DatabaseExists(conn, dbName)
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		// Удалим БД.
		if err = mysql.DatabaseDrop(conn, dbName); err != nil {
			log.Fatal(err)
		}
	}

	// Создадим БД.
	if err = mysql.DatabaseCreate(conn, dbName); err != nil {
		log.Fatal(err)
	}

	// Подключение к новой СУБД sample.
	connSample, err := db.Connect(&db.Config{
		Driver:   db.DriverMySQL,
		Host:     db.CfgDefaultHost,
		Port:     "3306",
		SSLMode:  db.SSLModeDisable,
		Database: "sample",
		User:     "root",
		Password: "secret",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Создадим таблицу foo.
	if _, err = connSample.Exec("CREATE TABLE foo(id varchar(255) primary key, name varchar(255))"); err != nil {
		log.Fatal(err)
	}

	// Создадим таблицу bar.
	if _, err = connSample.Exec("CREATE TABLE bar(id varchar(255) primary key, name varchar(255))"); err != nil {
		log.Fatal(err)
	}

	// Запросим таблицы, кроме bar.
	tables, err := mysql.TableList(connSample, "bar")
	if err != nil {
		log.Fatal(err)
	}

	// Удалим полученные таблицы.
	if err = mysql.TableDrop(connSample, tables...); err != nil {
		log.Fatal(err)
	}

	// Создадим представление bar_view.
	if _, err = connSample.Exec("CREATE VIEW bar_view(id, name) AS SELECT id, name FROM bar"); err != nil {
		log.Fatal(err)
	}

	// Запросим представления, кроме bar_view.
	views, err := mysql.ViewList(connSample, "bar_view")
	if err != nil {
		log.Fatal(err)
	}

	// Удалим полученные представления.
	if err = mysql.ViewDrop(connSample, views...); err != nil {
		log.Fatal(err)
	}
}
