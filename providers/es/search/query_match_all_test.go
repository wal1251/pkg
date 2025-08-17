package search_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/providers/es/search"
)

func TestQueryMatchAll_Empty(t *testing.T) {
	queryMatchAll := search.QueryMatchAll{}

	require.True(t, queryMatchAll.IsEmpty())
	require.Equal(t, search.QueryMatchAllKey, queryMatchAll.QueryType())

	jsonResponse := `{"match_all":{}}`

	res, err := queryMatchAll.MarshalJSON()
	require.NoError(t, err)

	require.JSONEq(t, jsonResponse, string(res))
}

func TestQueryMatchAll_WithBoost(t *testing.T) {
	var boost float32 = 1.0

	queryMatchAll := search.QueryMatchAll{
		Boost: &boost,
	}

	require.False(t, queryMatchAll.IsEmpty())
	require.Equal(t, search.QueryMatchAllKey, queryMatchAll.QueryType())

	jsonResponse := `{"match_all":{"boost":1}}`

	res, err := queryMatchAll.MarshalJSON()
	require.NoError(t, err)

	require.JSONEq(t, jsonResponse, string(res))
}
