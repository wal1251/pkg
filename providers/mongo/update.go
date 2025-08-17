package mongo

import "go.mongodb.org/mongo-driver/bson"

// Update оборачивает bson.D для удобного создания обновлений.
type Update struct {
	value bson.D
}

// NewUpdate создает новый объект обновления.
func NewUpdate() *Update {
	return &Update{value: bson.D{}}
}

// Set добавляет обновление с установкой поля (field = value).
func (u *Update) Set(field string, value interface{}) *Update {
	u.value = append(u.value, bson.E{Key: "$set", Value: bson.M{field: value}})

	return u
}

// Inc увеличивает значение поля на указанное значение (field += value).
func (u *Update) Inc(field string, value interface{}) *Update {
	u.value = append(u.value, bson.E{Key: "$inc", Value: bson.M{field: value}})

	return u
}

// Unset удаляет поле.
func (u *Update) Unset(field string) *Update {
	u.value = append(u.value, bson.E{Key: "$unset", Value: bson.M{field: ""}})

	return u
}

// Value возвращает bson.D для использования в запросах.
func (u *Update) Value() bson.D {
	return u.value
}
