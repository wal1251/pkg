package search

import (
	"fmt"

	"github.com/wal1251/pkg/providers/es/api"
	"github.com/wal1251/pkg/tools/serial"
)

const (
	QueryTermKey      = "term"
	QueryTermValueKey = "value"
	QueryTermBoostKey = "boost"
)

var _ Query = (*QueryTerm)(nil)

// QueryTerm запрос, который возвращает документы, содержащий точный term в указанном поле.
type QueryTerm struct {
	Field string
	Value string
	Boost *float32
}

func (q *QueryTerm) MarshalJSON() ([]byte, error) {
	return serial.ToBytes(api.Object{
		q.QueryType(): api.Object{
			q.Field: api.CreateObject(
				api.Property(QueryTermValueKey, q.Value),
				api.PropertyOmitempty(QueryTermBoostKey, q.Boost),
			),
		},
	}, serial.JSONEncode[api.Object])
}

func (q *QueryTerm) IsEmpty() bool {
	return q.Field == ""
}

func (q *QueryTerm) QueryType() string {
	return QueryTermKey
}

func NewQueryTerm(field string, value any) *QueryTerm {
	return &QueryTerm{
		Field: field,
		Value: fmt.Sprint(value),
	}
}
