package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObject_Merge(t *testing.T) {
	tests := []struct {
		name string
		o    Object
		arg  Object
		want Object
	}{
		{
			name: "Basic test",
			o: Object{
				"foo": Object{
					"bar": "1",
					"baz": "2",
				},
				"bar": Object{
					"foo": Array{1, 2, 3},
					"baz": "3",
				},
			},
			arg: Object{
				"foo1": "q",
				"foo": Object{
					"bar1": "3",
					"bar2": "4",
				},
				"bar": Object{
					"foo": Array{4, 5, 6, Object{"foo": 7}},
					"baz": 4,
				},
			},
			want: Object{
				"foo1": "q",
				"foo": Object{
					"bar":  "1",
					"baz":  "2",
					"bar1": "3",
					"bar2": "4",
				},
				"bar": Object{
					"foo": Array{1, 2, 3, 4, 5, 6, Object{"foo": 7}},
					"baz": 4,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.o.Merge(tt.arg))
		})
	}
}
