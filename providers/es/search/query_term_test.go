package search_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/providers/es/search"
)

func TestNewQueryTerm(t *testing.T) {
	queryTerm := search.NewQueryTerm("item", "item_value")

	require.False(t, queryTerm.IsEmpty())
	require.Equal(t, search.QueryTermKey, queryTerm.QueryType())

	jsonResponse := `{"term":{"item":{"value":"item_value"}}}`

	res, err := queryTerm.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, jsonResponse, string(res))
}
