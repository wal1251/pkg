package indices_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	indices2 "github.com/wal1251/pkg/providers/es/indices"
)

func TestFieldMapping(t *testing.T) {
	type outStruct struct {
		MappingType indices2.MappingType
		Property    string
	}

	typeKeyword := indices2.TypeKeyword
	typeText := indices2.TypeText

	analyzer := "russian"

	tests := []struct {
		name string
		in   reflect.StructField
		out  outStruct
	}{
		{
			name: "Status field",
			in: reflect.StructField{
				Name:      "Status",
				PkgPath:   "",
				Type:      reflect.TypeOf(""),
				Tag:       `es:"keyword" json:"Status"`,
				Offset:    16,
				Index:     []int{1},
				Anonymous: false,
			},
			out: outStruct{
				Property: "Status",
				MappingType: indices2.MappingType{
					Properties: nil,
					Type:       &typeKeyword,
					Analyzer:   nil,
				},
			},
		},
		{
			name: "ParentID field",
			in: reflect.StructField{
				Name:      "ParentID",
				PkgPath:   "",
				Type:      reflect.TypeOf(&uuid.Nil),
				Tag:       `es:"keyword" json:"ParentID,omitempty"`,
				Offset:    40,
				Index:     []int{3},
				Anonymous: false,
			},
			out: outStruct{
				Property: "ParentID",
				MappingType: indices2.MappingType{
					Properties: nil,
					Type:       &typeKeyword,
					Analyzer:   nil,
				},
			},
		},
		{
			name: "CaptionRu field",
			in: reflect.StructField{
				Name:      "CaptionRu",
				PkgPath:   "",
				Type:      reflect.TypeOf(""),
				Tag:       `search:"must,boost=2,lang=ru" es:"text,analyzer=russian" json:"CaptionRu"`,
				Offset:    48,
				Index:     []int{4},
				Anonymous: false,
			},
			out: outStruct{
				Property: "CaptionRu",
				MappingType: indices2.MappingType{
					Properties: nil,
					Type:       &typeText,
					Analyzer:   &analyzer,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mappingType, property, err := indices2.FieldMapping(tt.in)
			assert.NoError(t, err)
			assert.Equal(t, tt.out, outStruct{
				MappingType: mappingType,
				Property:    property,
			})
		})
	}
}
