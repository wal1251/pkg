package search

import (
	"fmt"

	"github.com/wal1251/pkg/providers/es/api"
	"github.com/wal1251/pkg/tools/serial"
)

const (
	QueryMultiMatchKey          = "multi_match"
	QueryMultiMatchQueryKey     = "query"
	QueryMultiMatchAnalyzerKey  = "analyzer"
	QueryMultiMatchFuzzinessKey = "fuzziness"
	QueryMultiMatchFieldsKey    = "fields"
)

var _ Query = (*QueryMultiMatch)(nil)

type (
	// QueryMultiMatch запрос для получения документов, основываясь на нескольких запросов сопоставления.
	// Запрос сопоставления(Match query) -- запрос для полнотекстового поиска, включая параметры нечеткого сопоставления.
	QueryMultiMatch struct {
		Query     string
		Analyzer  *string
		Fuzziness *string
		Fields    map[string]int
	}
)

func (q *QueryMultiMatch) MarshalJSON() ([]byte, error) {
	fields := make([]string, 0, len(q.Fields))

	for field, boost := range q.Fields {
		if boost > 1 {
			field = fmt.Sprintf("%s^%d", field, boost)
		}

		fields = append(fields, field)
	}

	return serial.ToBytes(api.Object{
		q.QueryType(): api.CreateObject(
			api.Property(QueryMultiMatchQueryKey, q.Query),
			api.PropertyOmitempty(QueryMultiMatchAnalyzerKey, q.Analyzer),
			api.PropertyOmitempty(QueryMultiMatchFuzzinessKey, q.Fuzziness),
			api.PropertyOmitempty(QueryMultiMatchFieldsKey, fields),
		),
	}, serial.JSONEncode[api.Object])
}

func (q *QueryMultiMatch) IsEmpty() bool {
	return q.Query == ""
}

func (q *QueryMultiMatch) QueryType() string {
	return QueryMultiMatchKey
}

func (q *QueryMultiMatch) AddFieldFromDescriptor(descriptors ...DocumentFieldDescriptor) {
	for _, descriptor := range descriptors {
		if q.Fields == nil {
			q.Fields = make(map[string]int)
		}

		boost := 1
		if descriptor.Boost != nil {
			boost = *descriptor.Boost
		}

		q.Fields[descriptor.Field] = boost
	}
}

func NewQueryMultiMatch(query string, fields ...DocumentFieldDescriptor) *QueryMultiMatch {
	multiMatch := &QueryMultiMatch{
		Query: query,
	}

	multiMatch.AddFieldFromDescriptor(fields...)

	return multiMatch
}
