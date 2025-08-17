package search

import (
	"github.com/wal1251/pkg/providers/es/api"
	"github.com/wal1251/pkg/tools/serial"
)

const (
	QueryExistsKey      = "exists"
	QueryExistsFieldKey = "field"
)

var _ Query = (*QueryExists)(nil)

// QueryExists запрос, который возвращает документы, содержащий индексированное значение для поля.
type QueryExists struct {
	Field string
}

func (q *QueryExists) MarshalJSON() ([]byte, error) {
	return serial.ToBytes(api.Object{
		q.QueryType(): api.Object{
			QueryExistsFieldKey: q.Field,
		},
	}, serial.JSONEncode[api.Object])
}

func (q *QueryExists) IsEmpty() bool {
	return q.Field == ""
}

func (q *QueryExists) QueryType() string {
	return QueryExistsKey
}
