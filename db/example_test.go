package db_test

import (
	"context"
	"log"

	"github.com/wal1251/pkg/db"
)

func ExampleMigrate() {
	// Создаем подключение к СУБД.
	cfg := db.NewCfgSQLiteMem("db")

	// Источник миграции должен быть задан.
	cfg.Migration = "./migrate"
	cfg.Debug = true

	// Мигрируем.
	if err := db.Migrate(context.TODO(), cfg, db.MigrateDefault(cfg)); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
}
