package search_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/providers/es/search"
)

func TestQueryExists(t *testing.T) {
	queryExists := search.QueryExists{
		Field: "item",
	}

	require.False(t, queryExists.IsEmpty())
	require.Equal(t, search.QueryExistsKey, queryExists.QueryType())

	responseJson := `{"exists":{"field":"item"}}`

	res, err := queryExists.MarshalJSON()
	require.NoError(t, err)

	require.JSONEq(t, responseJson, string(res))
}
