package collections_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/tools/collections"
)

func TestNewSet(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{
			name:   "Basic test",
			values: []int{1, 2, 3, 4},
			want:   []int{1, 2, 3, 4},
		},
		{
			name:   "Duplicates test",
			values: []int{1, 2, 3, 3, 3, 4},
			want:   []int{1, 2, 3, 4},
		},
		{
			name:   "Empty test",
			values: []int{},
			want:   []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := collections.NewSet(tt.values...)
			assert.ElementsMatch(t, tt.want, set.ToSlice())
		})
	}
}

func TestSet_Add(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{
			name:   "Basic test",
			values: []int{1, 2, 3, 4},
			want:   []int{1, 2, 3, 4},
		},
		{
			name:   "Duplicates test",
			values: []int{1, 2, 3, 3, 3, 4},
			want:   []int{1, 2, 3, 4},
		},
		{
			name:   "Empty test",
			values: []int{},
			want:   []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := collections.NewSet[int]()
			set.Add(tt.values...)
			assert.ElementsMatch(t, tt.want, set.ToSlice())
		})
	}
}

func TestSet_Contains(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		arg    int
		want   bool
	}{
		{
			name:   "Contains value",
			values: []int{1, 2, 3, 4},
			arg:    3,
			want:   true,
		},
		{
			name:   "Not contains value",
			values: []int{1, 2, 3, 4},
			arg:    5,
		},
		{
			name:   "Contains value with duplicates",
			values: []int{1, 2, 3, 3, 3, 4},
			arg:    3,
			want:   true,
		},
		{
			name:   "Empty test",
			values: []int{},
			arg:    3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := collections.NewSet(tt.values...)
			assert.Equal(t, tt.want, set.Contains(tt.arg))
		})
	}
}

func TestSet_ForEach(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{
			name:   "Basic test",
			values: []int{1, 2, 3, 4},
			want:   []int{1, 2, 3, 4},
		},
		{
			name:   "Duplicates test",
			values: []int{1, 2, 3, 3, 3, 4},
			want:   []int{1, 2, 3, 4},
		},
		{
			name:   "Empty test",
			values: []int{},
			want:   []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]int, 0)
			set := collections.NewSet(tt.values...)
			set.ForEach(func(i int) {
				result = append(result, i)
			})
			assert.ElementsMatch(t, tt.want, result)
		})
	}
}

func TestSet_Len(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   int
	}{
		{
			name:   "Basic test",
			values: []int{1, 2, 3, 4},
			want:   4,
		},
		{
			name:   "Duplicates test",
			values: []int{1, 2, 3, 3, 3, 4},
			want:   4,
		},
		{
			name:   "Empty test",
			values: []int{},
			want:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, collections.NewSet(tt.values...).Len())
		})
	}
}

func TestSet_Remove(t *testing.T) {
	tests := []struct {
		name      string
		values    []int
		arg       int
		wantSlice []int
	}{
		{
			name:      "Remove value",
			values:    []int{1, 2, 3, 4},
			arg:       3,
			wantSlice: []int{1, 2, 4},
		},
		{
			name:      "Remove not contained value",
			values:    []int{1, 2, 3, 4},
			arg:       5,
			wantSlice: []int{1, 2, 3, 4},
		},
		{
			name:      "Remove value with duplicates",
			values:    []int{1, 2, 3, 3, 3, 4},
			arg:       3,
			wantSlice: []int{1, 2, 4},
		},
		{
			name:      "Empty test",
			values:    []int{},
			arg:       3,
			wantSlice: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := collections.NewSet(tt.values...)
			set.Remove(tt.arg)
			assert.ElementsMatch(t, tt.wantSlice, set.ToSlice())
		})
	}
}

func TestSet_ToSlice(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{
			name:   "Basic test",
			values: []int{1, 2, 3, 4},
			want:   []int{1, 2, 3, 4},
		},
		{
			name:   "Duplicates test",
			values: []int{1, 2, 3, 3, 3, 4},
			want:   []int{1, 2, 3, 4},
		},
		{
			name:   "Empty test",
			values: []int{},
			want:   []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := collections.NewSet(tt.values...)
			assert.ElementsMatch(t, tt.want, set.ToSlice())
		})
	}
}
