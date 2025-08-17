package search_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/providers/es/search"
)

func TestQueryDescriptor(t *testing.T) {
	tests := []struct {
		name string
		in   any
		out  search.QueryDescriptor
	}{
		{
			name: "empty struct",
			in:   struct{}{},
			out:  make(search.QueryDescriptor, 0),
		},
		{
			name: "struct without tag",
			in: struct {
				Name string
				Size int
			}{
				Name: "file",
				Size: 0,
			},
			out: make(search.QueryDescriptor, 0),
		},
		{
			name: "struct with search tag -- query",
			in: struct {
				Name string `search:"query"`
			}{
				Name: "file",
			},
			out: []search.QueryFieldDescriptor{
				{
					Name:  "",
					Type:  "query",
					Field: "Name",
				},
			},
		},
		{
			name: "struct with search tag -- filter",
			in: struct {
				Name string `search:"filter"`
			}{
				Name: "file",
			},
			out: []search.QueryFieldDescriptor{
				{
					Name:  "",
					Type:  "filter",
					Field: "Name",
				},
			},
		},
		{
			name: "struct with search tags -- filter",
			in: struct {
				Name string `search:"filter"`
			}{
				Name: "file",
			},
			out: []search.QueryFieldDescriptor{
				{
					Name:  "",
					Type:  "filter",
					Field: "Name",
				},
			},
		},
		{
			name: "struct with few fields with search tag -- filter, query",
			in: struct {
				AuthorID  uuid.UUID `search:"query"`
				CenturyID uuid.UUID `search:"filter"`
				YearID    uuid.UUID `search:"query"`
				ArtFormID uuid.UUID `search:"filter"`
			}{
				AuthorID:  uuid.Nil,
				CenturyID: uuid.Nil,
				YearID:    uuid.Nil,
				ArtFormID: uuid.Nil,
			},
			out: []search.QueryFieldDescriptor{
				{
					Name:  "",
					Type:  "query",
					Field: "AuthorID",
				},
				{
					Name:  "",
					Type:  "filter",
					Field: "CenturyID",
				},
				{
					Name:  "",
					Type:  "query",
					Field: "YearID",
				},
				{
					Name:  "",
					Type:  "filter",
					Field: "ArtFormID",
				},
			},
		},
		{
			name: "struct with search type tag",
			in: struct {
				Name string `search:"type=1, type=2"`
			}{
				Name: "file",
			},
			out: []search.QueryFieldDescriptor{
				{
					Name:  "",
					Type:  "2",
					Field: "Name",
				},
			},
		},
		{
			name: "struct with search name tag",
			in: struct {
				Name string `search:"type=type, name=name"`
			}{
				Name: "file",
			},
			out: []search.QueryFieldDescriptor{
				{
					Name:  "name",
					Type:  "type",
					Field: "Name",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := search.QueryDescribe(tt.in)
			assert.NoError(t, err)
			assert.Equal(t, tt.out, out)
		})
	}
}

func TestDocumentFieldsDescription(t *testing.T) {
	type Folder struct {
		ID             uuid.UUID  `es:"keyword" json:"ID"`
		Status         string     `es:"keyword" json:"Status"`
		IsArchived     bool       `es:"boolean" json:"IsArchived"`
		ParentID       *uuid.UUID `es:"keyword" json:"ParentID,omitempty"`
		CaptionRu      string     `search:"must,boost=2,lang=ru" es:"text,analyzer=russian" json:"CaptionRu"`
		CaptionEn      string     `search:"must,boost=2,lang=en" es:"text,analyzer=english" json:"CaptionEn"`
		DescriptionRu  *string    `search:"must,lang=ru" es:"text,analyzer=russian" json:"DescriptionRu,omitempty"`
		DescriptionEn  *string    `search:"must,lang=en" es:"text,analyzer=english" json:"DescriptionEn,omitempty"`
		MaterialsKinds []string   `es:"keyword" json:"MaterialsKinds,omitempty"`
		Tags           []string   `es:"keyword" json:"Tags,omitempty"`
	}

	out, err := search.DocumentFieldsDescription(reflect.TypeOf((any)((*Folder)(nil))).Elem())
	assert.NoError(t, err)
	boost := 2
	langEn := "en"
	langRu := "ru"
	expected := []search.DocumentFieldDescriptor{
		{
			Name:     "",
			Type:     "filter",
			Field:    "ID",
			Boost:    nil,
			Language: nil,
		},
		{
			Name:     "",
			Type:     "filter",
			Field:    "Status",
			Boost:    nil,
			Language: nil,
		},
		{
			Name:     "",
			Type:     "filter",
			Field:    "IsArchived",
			Boost:    nil,
			Language: nil,
		},
		{
			Name:     "",
			Type:     "filter",
			Field:    "ParentID",
			Boost:    nil,
			Language: nil,
		},
		{
			Name:     "",
			Type:     "must",
			Field:    "CaptionRu",
			Boost:    &boost,
			Language: &langRu,
		},
		{
			Name:     "",
			Type:     "must",
			Field:    "CaptionEn",
			Boost:    &boost,
			Language: &langEn,
		},
		{
			Name:     "",
			Type:     "must",
			Field:    "DescriptionRu",
			Boost:    nil,
			Language: &langRu,
		},
		{
			Name:     "",
			Type:     "must",
			Field:    "DescriptionEn",
			Boost:    nil,
			Language: &langEn,
		},
		{
			Name:     "",
			Type:     "filter",
			Field:    "MaterialsKinds",
			Boost:    nil,
			Language: nil,
		},
		{
			Name:     "",
			Type:     "filter",
			Field:    "Tags",
			Boost:    nil,
			Language: nil,
		},
	}

	require.Equal(t, expected, out)
}
