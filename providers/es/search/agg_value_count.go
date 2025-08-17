package search

import (
	"github.com/wal1251/pkg/providers/es/api"
	"github.com/wal1251/pkg/tools/serial"
)

const (
	AggregationValueCountKey      = "value_count"
	AggregationValueCountFieldKey = "field"
)

var _ Aggregation = (*AggregationValueCount)(nil)

type (
	// AggregationValueCount запрос подсчета количества значений, извлекаемых из агрегированных документов.
	AggregationValueCount struct {
		Field string
	}
)

func (a *AggregationValueCount) MarshalJSON() ([]byte, error) {
	return serial.ToBytes(api.Object{
		a.AggregationType(): api.Object{
			AggregationValueCountFieldKey: a.Field,
		},
	}, serial.JSONEncode[api.Object])
}

func (a *AggregationValueCount) AggregationType() string {
	return AggregationValueCountKey
}
