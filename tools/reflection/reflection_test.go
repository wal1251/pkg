package reflection_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/tools/reflection"
)

func TestTagParseKeyValue(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out1 string
		out2 map[string]string
	}{
		{
			name: "query, empty",
			in:   "",
			out1: "",
			out2: map[string]string{},
		},
		{
			name: "query",
			in:   "query",
			out1: "query",
			out2: map[string]string{},
		},
		{
			name: "query, t1=1, t2=3",
			in:   "query, t1=1, t2=3",
			out1: "query",
			out2: map[string]string{
				"t1": "1",
				"t2": "3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out1, out2 := reflection.TagParseKeyValue(tt.in)
			require.Equal(t, tt.out1, out1)
			require.Equal(t, tt.out2, out2)
		})
	}
}

func TestGetJSONName(t *testing.T) {
	type (
		j struct {
			Field1 string
			Field2 string `json:"field_2"`
			Field3 string `json:"-,field_3"`
		}

		out struct {
			out1 string
			out2 bool
		}
	)

	cases := []struct {
		name string
		in   reflect.StructField
		out  out
	}{
		{
			name: "Имя в json совпадает с именем поля в структуре",
			in:   reflect.StructField{Name: "Field1"},
			out:  out{out1: "Field1", out2: true},
		},
		{
			name: "Имя в json совпадает с первым значением в tag",
			in:   reflect.StructField{Name: "Field2", Tag: `json:"field_2"`},
			out:  out{out1: "field_2", out2: true},
		},
		{
			name: "Имя в json отсутствует у данного поля",
			in:   reflect.StructField{},
			out:  out{out1: "", out2: true},
		},
		{
			name: "Имя в json tag содержит прочерк",
			in:   reflect.StructField{Tag: `json:"-,tag1"`},
			out:  out{out1: "", out2: false},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			out1, out2 := reflection.GetJSONName(tt.in)
			assert.Equal(t, tt.out.out1, out1)
			assert.Equal(t, tt.out.out2, out2)
		})
	}
}

func TestMapByTypeName(t *testing.T) {
	type (
		LocalType1 string
		LocalType2 string
		LocalType3 string
	)

	tests := []struct {
		name string
		in   []any
		out  map[string]any
	}{
		{
			name: "[int, string float64]",
			in:   []any{1, "abc", 1.1},
			out: map[string]any{
				"int":     1,
				"string":  "abc",
				"float64": 1.1,
			},
		},
		{
			name: "[LocalType1, LocalType2, LocalType3]",
			in:   []any{LocalType1("1"), LocalType2("2"), LocalType3("3")},
			out: map[string]any{
				"LocalType1": LocalType1("1"),
				"LocalType2": LocalType2("2"),
				"LocalType3": LocalType3("3"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.out, reflection.MapByTypeName(tt.in...))
		})
	}
}

func TestIsEmptyValue(t *testing.T) {
	type testData struct {
		value    interface{}
		expected bool
	}

	testCases := []testData{
		{value: "", expected: true},
		{value: "hello", expected: false},
		{value: []int{}, expected: true},
		{value: []int{1, 2, 3}, expected: false},
		{value: map[string]int{}, expected: true},
		{value: map[string]int{"a": 1}, expected: false},
		{value: 0, expected: true},
		{value: 1, expected: false},
		{value: false, expected: true},
		{value: true, expected: false},
	}

	for _, testCase := range testCases {
		result := reflection.IsEmptyValue(testCase.value)
		if result != testCase.expected {
			t.Errorf("IsEmptyValue(%v) = %v; want %v", testCase.value, result, testCase.expected)
		}
	}
}
