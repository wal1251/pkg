package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Document представляет собой документ MongoDB.
// Это карта, где ключи — названия полей, а значения — их содержимое.
type Document map[string]interface{}

// GetID возвращает ID документа.
func (d Document) GetID() (ID, error) {
	rawID, ok := d["_id"]
	if !ok {
		return ID{}, ErrIDNotFoundInDocument
	}

	objID, ok := rawID.(primitive.ObjectID)
	if !ok {
		return ID{}, ErrInvalidIDInDocument
	}

	return ID{value: objID}, nil
}

func (d Document) GetStringValue(key string) (string, error) {
	value, ok := d[key]
	if !ok {
		return "", ErrDocumentValueNotFound
	}

	str, ok := value.(string)
	if !ok {
		return "", ErrStringValueConversionFailed
	}

	return str, nil
}

// HasKey проверяет наличие ключа в документе.
func (d Document) HasKey(key string) bool {
	_, exists := d[key]

	return exists
}

// GetKeys возвращает список ключей документа.
func (d Document) GetKeys() []string {
	keys := make([]string, 0, len(d))
	for key := range d {
		keys = append(keys, key)
	}

	return keys
}

// GetValue возвращает значение документа по ключу.
func (d Document) GetValue(key string) (interface{}, error) {
	value, ok := d[key]
	if !ok {
		return nil, ErrDocumentValueNotFound
	}

	return value, nil
}
