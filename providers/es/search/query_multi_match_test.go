package search_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	search2 "github.com/wal1251/pkg/providers/es/search"
)

func TestQueryMultiMatchWithoutBoost(t *testing.T) {
	queryMultiMatch := search2.NewQueryMultiMatch(
		"item",
		search2.DocumentFieldDescriptor{
			Name:  "document_name_3",
			Type:  "document_type_3",
			Field: "document_field_3",
		},
	)

	jsonResponse := `{"multi_match":{"fields":["document_field_3"],"query":"item"}}`

	require.False(t, queryMultiMatch.IsEmpty())
	require.Equal(t, search2.QueryMultiMatchKey, queryMultiMatch.QueryType())

	res, err := queryMultiMatch.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, jsonResponse, string(res))
}

func TestQueryMultiMatchWithBoost(t *testing.T) {
	boost := 2
	language := "ru"
	queryMultiMatch := search2.NewQueryMultiMatch(
		"item",
		search2.DocumentFieldDescriptor{
			Name:     "document_name_3",
			Type:     "document_type_3",
			Field:    "document_field_3",
			Boost:    &boost,
			Language: &language,
		},
	)

	jsonResponse := `{"multi_match":{"fields":["document_field_3^2"],"query":"item"}}`

	require.False(t, queryMultiMatch.IsEmpty())
	require.Equal(t, search2.QueryMultiMatchKey, queryMultiMatch.QueryType())

	res, err := queryMultiMatch.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, jsonResponse, string(res))
}
