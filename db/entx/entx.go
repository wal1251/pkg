// Package entx расширяет возможности библиотеки ent.
package entx

import (
	"entgo.io/ent/dialect/sql"

	"github.com/wal1251/pkg/db"
)

const (
	FieldID          = "id"           // Идентификатор - первичный ключ записи.
	FieldIsDeleted   = "is_deleted"   // Флаг удаляемой записи.
	FieldDeletedTime = "deleted_time" // Время установки флага удаления записи.
	FieldCreateTime  = "create_time"  // Время создания записи.
	FieldUpdateTime  = "update_time"  // Время обновления записи.
)

// Driver возвращает ent драйвер для переданных параметров подключения к БД.
func Driver(cfg db.ConnectionDescriber) (*sql.Driver, error) {
	conn, err := db.Connect(cfg)
	if err != nil {
		return nil, err
	}

	return sql.OpenDB(cfg.DriverName(), conn), nil
}
