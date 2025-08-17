package search_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/providers/es/search"
)

func TestQueryBooleanWithoutDocumentFieldDescriptorList(t *testing.T) {
	type Search struct {
		Language            string
		Query               string      `search:"query"`
		IsArchived          *bool       `search:"filter"`
		AuthorID            *uuid.UUID  `search:"filter"`
		CenturyID           *uuid.UUID  `search:"filter"`
		YearID              *uuid.UUID  `search:"filter"`
		ArtFormID           *uuid.UUID  `search:"filter"`
		MaterialsTechnicsID *uuid.UUID  `search:"filter"`
		CenturyHalf         *int        `search:"filter"`
		Tags                []uuid.UUID `search:"filter"`
	}

	searchTest := Search{Language: "ru", Query: "item"}
	descriptor, err := search.QueryDescribe(searchTest)
	require.NoError(t, err)

	queryBoolean := search.NewQueryBoolean(
		searchTest,
		searchTest.Language,
		descriptor,
		make([]search.DocumentFieldDescriptor, 0),
	)

	require.True(t, queryBoolean.IsEmpty())

	jsonResponse := `{"bool":{}}`
	res, err := queryBoolean.MarshalJSON()
	require.NoError(t, err)

	require.JSONEq(t, jsonResponse, string(res))
}

// QueryBoolean with DocumentFieldDescriptor list
func TestQueryBooleanWithDocumentFieldDescriptorList(t *testing.T) {
	type Search struct {
		Language            string
		Query               string      `search:"query"`
		IsArchived          *bool       `search:"filter"`
		AuthorID            *uuid.UUID  `search:"filter"`
		CenturyID           *uuid.UUID  `search:"filter"`
		YearID              *uuid.UUID  `search:"filter"`
		ArtFormID           *uuid.UUID  `search:"filter"`
		MaterialsTechnicsID *uuid.UUID  `search:"filter"`
		CenturyHalf         *int        `search:"filter"`
		Tags                []uuid.UUID `search:"filter"`
	}

	searchTest := Search{Language: "ru", Query: "item"}

	type Folder struct {
		ID             uuid.UUID  `es:"keyword" json:"ID"`
		Status         string     `es:"keyword" json:"Status"`
		IsArchived     bool       `es:"boolean" json:"IsArchived"`
		ParentID       *uuid.UUID `es:"keyword" json:"ParentID,omitempty"`
		CaptionRu      string     `search:"must,boost=2,lang=ru" es:"text,analyzer=russian" json:"CaptionRu"`
		CaptionEn      string     `search:"must,boost=2,lang=en" es:"text,analyzer=english" json:"CaptionEn"`
		DescriptionEn  *string    `search:"must,lang=en" es:"text,analyzer=english" json:"DescriptionEn,omitempty"`
		MaterialsKinds []string   `es:"keyword" json:"MaterialsKinds,omitempty"`
		Tags           []string   `es:"keyword" json:"Tags,omitempty"`
	}

	descriptor, err := search.QueryDescribe(searchTest)
	require.NoError(t, err)

	documentFieldsDescription, err := search.DocumentFieldsDescription(reflect.TypeOf((any)((*Folder)(nil))).Elem())
	require.NoError(t, err)

	queryBoolean := search.NewQueryBoolean(
		searchTest,
		searchTest.Language,
		descriptor,
		documentFieldsDescription,
	)

	require.False(t, queryBoolean.IsEmpty())

	jsonResponse := `{"bool":{"must":[{"multi_match":{"fields":["CaptionRu^2"],"query":"item"}}]}}`
	res, err := queryBoolean.MarshalJSON()
	require.NoError(t, err)

	require.JSONEq(t, jsonResponse, string(res))
}
