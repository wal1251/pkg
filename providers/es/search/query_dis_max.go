package search

import (
	"github.com/wal1251/pkg/providers/es/api"
	"github.com/wal1251/pkg/tools/serial"
)

var _ Query = (*QueryDisjunctionMax)(nil)

const (
	QueryDisMax           = "dis_max"
	QueryDisMaxQueriesKey = "queries"
	QueryDisMaxTieBreaker = "tie_breaker"
)

// QueryDisjunctionMax запрос на дизъюнкцию логических предложений.
// Возвращаемые документы соответствуют с одном или несколькими логическими предложениями.
type QueryDisjunctionMax struct {
	Queries    Queries
	TieBreaker float32
}

func (q QueryDisjunctionMax) MarshalJSON() ([]byte, error) {
	return serial.ToBytes(api.Object{
		q.QueryType(): api.CreateObject(
			api.Property(QueryDisMaxQueriesKey, q.Queries),
			api.PropertyOmitempty(QueryDisMaxTieBreaker, q.TieBreaker),
		),
	}, serial.JSONEncode[api.Object])
}

func (q QueryDisjunctionMax) QueryType() string {
	return QueryDisMax
}

func (q QueryDisjunctionMax) IsEmpty() bool {
	return len(q.Queries) == 0
}

func NewQueryDisjunctionMax(queries Queries) *QueryDisjunctionMax {
	return &QueryDisjunctionMax{
		Queries: queries,
	}
}
