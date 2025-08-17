package search_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/providers/es/search"
)

func TestAggregationValueCount(t *testing.T) {
	aggregationValueCount := search.AggregationValueCount{
		Field: "item",
	}

	require.Equal(t, search.AggregationValueCountKey, aggregationValueCount.AggregationType())

	jsonResponse := `{"value_count":{"field":"item"}}`

	res, err := aggregationValueCount.MarshalJSON()
	require.NoError(t, err)

	require.JSONEq(t, jsonResponse, string(res))
}
