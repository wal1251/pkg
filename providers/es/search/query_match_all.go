package search

import (
	"github.com/wal1251/pkg/providers/es/api"
	"github.com/wal1251/pkg/tools/serial"
)

const (
	QueryMatchAllKey      = "match_all"
	QueryMatchAllBoostKey = "boost"
)

var _ Query = (*QueryMatchAll)(nil)

// QueryMatchAll запрос, который сопоставляет все документы, давая им всем _score 1.0.
// Парамер _score можно изменить параметром boost.
type QueryMatchAll struct {
	Boost *float32
}

func (q *QueryMatchAll) MarshalJSON() ([]byte, error) {
	return serial.ToBytes(api.Object{
		q.QueryType(): api.CreateObject(
			api.PropertyOmitempty(QueryMatchAllBoostKey, q.Boost),
		),
	}, serial.JSONEncode[api.Object])
}

func (q *QueryMatchAll) IsEmpty() bool {
	return q.Boost == nil
}

func (q *QueryMatchAll) QueryType() string {
	return QueryMatchAllKey
}
