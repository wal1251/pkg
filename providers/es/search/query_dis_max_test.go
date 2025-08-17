package search_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	search2 "github.com/wal1251/pkg/providers/es/search"
)

func TestQueryDisjunctionMax(t *testing.T) {
	queryDisjunctionMax := search2.NewQueryDisjunctionMax(
		search2.Queries{
			search2.NewQueryTerm("item_1", "item_1_value"),
			search2.NewQueryTerm("item_2", "item_2_value"),
			search2.NewQueryTerm("item_3", "item_3_value"),
		})

	require.False(t, queryDisjunctionMax.IsEmpty())
	require.Equal(t, search2.QueryDisMax, queryDisjunctionMax.QueryType())

	jsonResponse := `{"dis_max":{"queries":[{"term":{"item_1":{"value":"item_1_value"}}},{"term":{"item_2":{"value":"item_2_value"}}},{"term": {"item_3": {"value":"item_3_value"}}}]}}`

	res, err := queryDisjunctionMax.MarshalJSON()
	require.NoError(t, err)

	require.JSONEq(t, jsonResponse, string(res))
}
