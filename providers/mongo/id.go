package mongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ID struct {
		value primitive.ObjectID
	}
)

// NewID создает новый ID.
func NewID() ID {
	return ID{value: primitive.NewObjectID()}
}

// ParseID парсит строку в ID.
func ParseID(id string) (ID, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ID{}, fmt.Errorf("failed to parse id: %w", err)
	}

	return ID{value: objID}, nil
}

// Hex возвращает строковое представление ID.
func (id ID) Hex() string {
	return id.value.Hex()
}

// Value возвращает внутренний ObjectID (для использования внутри пакета).
func (id ID) Value() primitive.ObjectID {
	return id.value
}
