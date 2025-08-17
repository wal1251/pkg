package collections_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/tools/collections"
)

func TestSlice_Chunked(t *testing.T) {
	seq := func(n int) []int {
		s := make([]int, n)
		for i := 0; i < n; i++ {
			s[i] = i + 1
		}
		return s
	}

	tests := []struct {
		name           string
		s              []int
		size           int
		wantChunkSizes []int
		wantSize       int
	}{
		{
			name:           "Базовый кейс",
			s:              seq(10),
			size:           3,
			wantChunkSizes: []int{3, 3, 3, 1},
			wantSize:       4,
		},
		{
			name:           "Пустой исходный слайс",
			s:              seq(0),
			size:           3,
			wantChunkSizes: []int{},
			wantSize:       0,
		},
		{
			name:           "Размер куска превышает исходный слайс",
			s:              seq(10),
			size:           13,
			wantChunkSizes: []int{10},
			wantSize:       1,
		},
		{
			name:           "Размер куска = 0",
			s:              seq(0),
			size:           0,
			wantChunkSizes: []int{},
			wantSize:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := collections.Chunked(tt.s, tt.size)
			if assert.Equal(t, tt.wantSize, len(result)) {
				actualSizes := make([]int, len(tt.wantChunkSizes))
				for i, chunk := range result {
					actualSizes[i] = len(chunk)
				}
				assert.Equal(t, tt.wantChunkSizes, actualSizes)
			}
		})
	}
}

func TestSlice_Map_Int(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		fn   func(int) int
		out  []int
	}{
		{
			name: "Пустой вход",
			in:   []int{},
			fn:   func(i int) int { return i },
			out:  []int{},
		},
		{
			name: "Умножение на 2",
			in:   []int{1, 2, 3, 4, 5, 6, 7, 8},
			fn:   func(i int) int { return i * 2 },
			out:  []int{2, 4, 6, 8, 10, 12, 14, 16},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.out, collections.Map(tt.in, tt.fn))
		})
	}
}

func TestSlice_Map_IntToString(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		fn   func(int) string
		out  []string
	}{
		{
			name: "Пустой вход",
			in:   []int{},
			fn:   func(i int) string { return fmt.Sprintf("%d", i) },
			out:  []string{},
		},
		{
			name: "Конвертация в string",
			in:   []int{1, 2, 3, 4, 5, 6, 7, 8},
			fn:   func(i int) string { return fmt.Sprintf("%d", i) },
			out:  []string{"1", "2", "3", "4", "5", "6", "7", "8"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.out, collections.Map(tt.in, tt.fn))
		})
	}
}

func TestSlice_FlatMap(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		fn   func(int) []int
		out  []int
	}{
		{
			name: "Пустой вход",
			in:   []int{},
			fn:   func(i int) []int { return []int{} },
			out:  []int{},
		},
		{
			name: "Дубликаты",
			in:   []int{1, 2, 3, 4},
			fn:   func(i int) []int { return []int{i, i} },
			out:  []int{1, 1, 2, 2, 3, 3, 4, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.out, collections.FlatMap(tt.in, tt.fn))
		})
	}
}

func TestSlice_Filter(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		fn   func(int) bool
		out  []int
	}{
		{
			name: "Пустой вход",
			in:   []int{},
			fn:   func(i int) bool { return true },
			out:  []int{},
		},
		{
			name: "Четность",
			in:   []int{1, 2, 3, 4, 5, 6, 7, 8},
			fn:   func(i int) bool { return i%2 == 0 },
			out:  []int{2, 4, 6, 8},
		},
		{
			name: "Всегда False",
			in:   []int{1, 2, 3, 4, 5, 6, 7, 8},
			fn:   func(i int) bool { return false },
			out:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.out, collections.Filter(tt.in, tt.fn))
		})
	}
}

func TestSlice_Split(t *testing.T) {
	type Out struct {
		first  []int
		second []int
	}

	tests := []struct {
		name string
		in   []int
		fn   func(int) bool
		out  Out
	}{
		{
			name: "Пустой вход",
			in:   []int{},
			fn:   func(i int) bool { return true },
			out:  Out{[]int{}, []int{}},
		},
		{
			name: "Четность и нечетность",
			in:   []int{1, 2, 3, 4, 5, 6, 7, 8},
			fn:   func(i int) bool { return i%2 == 0 },
			out:  Out{[]int{2, 4, 6, 8}, []int{1, 3, 5, 7}},
		},
		{
			name: "Всегда False",
			in:   []int{1, 2, 3, 4, 5, 6, 7, 8},
			fn:   func(i int) bool { return false },
			out:  Out{[]int{}, []int{1, 2, 3, 4, 5, 6, 7, 8}},
		},
		{
			name: "Всегда True",
			in:   []int{1, 2, 3, 4, 5, 6, 7, 8},
			fn:   func(i int) bool { return true },
			out:  Out{[]int{1, 2, 3, 4, 5, 6, 7, 8}, []int{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			firstOut, secondOut := collections.Split(tt.in, tt.fn)
			assert.ElementsMatch(t, tt.out.first, firstOut)
			assert.ElementsMatch(t, tt.out.second, secondOut)
		})
	}
}

func TestSlice_MapWithErr(t *testing.T) {
	type Out struct {
		list []int
		err  error
	}

	tests := []struct {
		name string
		in   []int
		fn   func(int) (int, error)
		out  Out
	}{
		{
			name: "Пустой вход",
			in:   []int{},
			fn:   func(i int) (int, error) { return i, nil },
			out:  Out{[]int{}, nil},
		},
		{
			name: "Некорректный второй элемент",
			in:   []int{1, 2},
			fn: func(i int) (int, error) {
				if i%2 == 0 {
					return 0, errors.New("четный элемент")
				}
				return i, nil
			},
			out: Out{nil, errors.New("четный элемент")},
		},
		{
			name: "Без ошибок",
			in:   []int{1, 2, 3, 4, 5},
			fn:   func(i int) (int, error) { return i, nil },
			out:  Out{[]int{1, 2, 3, 4, 5}, nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := collections.MapWithErr(tt.in, tt.fn)
			assert.Equal(t, tt.out.err, err)
			assert.ElementsMatch(t, tt.out.list, list)
		})
	}
}

func TestSlice_MapChunkedWithErr(t *testing.T) {
	type (
		In struct {
			list []int
			size int
		}

		Out struct {
			list []int
			err  error
		}
	)
	tests := []struct {
		name string
		in   In
		fn   func([]int) ([]int, error)
		out  Out
	}{
		{
			name: "Пустой вход",
			in: In{
				list: []int{},
				size: 0,
			},
			fn:  func(l []int) ([]int, error) { return l, nil },
			out: Out{[]int{}, nil},
		},
		{
			name: "len < size",
			in: In{
				list: []int{1, 2, 3, 4, 5},
				size: 10,
			},
			fn:  func(l []int) ([]int, error) { return l, nil },
			out: Out{[]int{1, 2, 3, 4, 5}, nil},
		},
		{
			name: "len > size",
			in: In{
				list: []int{1, 2, 3, 4, 5},
				size: 1,
			},
			fn:  func(l []int) ([]int, error) { return l, nil },
			out: Out{[]int{1, 2, 3, 4, 5}, nil},
		},
		{
			name: "Ошибка в разбивке",
			in: In{
				list: []int{1, 2, 3, 4, 5},
				size: 1,
			},
			fn:  func(l []int) ([]int, error) { return l, errors.New("err") },
			out: Out{[]int{}, errors.New("err")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := collections.MapChunkedWithErr(tt.in.list, tt.in.size, tt.fn)
			assert.Equal(t, tt.out.err, err)
			assert.ElementsMatch(t, tt.out.list, list)
		})
	}
}

func TestSlice_NonNil(t *testing.T) {
	getStrPointer := func(str string) *string {
		return &str
	}
	strPointer1, strPointer2, strPointer3 := getStrPointer("1"), getStrPointer("2"), getStrPointer("3")
	tests := []struct {
		name string
		in   []*string
		out  []string
	}{
		{
			name: "Пустой вход",
			in:   []*string{},
			out:  []string{},
		},
		{
			name: "Смешанный слайс",
			in:   []*string{strPointer1, nil, strPointer2, nil, strPointer3, nil},
			out:  []string{*strPointer1, *strPointer2, *strPointer3},
		},
		{
			name: "Слайс из Nil",
			in:   []*string{nil, nil, nil},
			out:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.out, collections.NonNil(tt.in))
		})
	}
}

func TestSlice_ExceptIndexes(t *testing.T) {
	type In struct {
		list []int
		set  collections.Set[int]
	}

	tests := []struct {
		name string
		in   In
		out  []int
	}{
		{
			name: "Пустой вход",
			in: In{
				list: []int{},
				set:  collections.Set[int]{},
			},
			out: []int{},
		},
		{
			name: "Частичная вхождение индексов в Set",
			in: In{
				list: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				set:  collections.NewSet(1, 5, 8, 9, 10),
			},
			out: []int{0, 2, 3, 4, 6, 7},
		},
		{
			name: "Полное вхождение индексов в Set",
			in: In{
				list: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				set:  collections.NewSet(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
			},
			out: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.out, collections.ExceptIndexes(tt.in.list, tt.in.set))
		})
	}
}

func TestSlice_KeysOfMap(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]int
		out  []string
	}{
		{
			name: "Пустой вход",
			in:   make(map[string]int),
			out:  []string{},
		},
		{
			name: "Заполненная map[string]int",
			in: map[string]int{
				"1": 1,
				"2": 2,
				"3": 3,
			},
			out: []string{"1", "2", "3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.out, collections.KeysOfMap(tt.in))
		})
	}
}

func TestSlice_ValuesOfMap(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]int
		out  []int
	}{
		{
			name: "Пустой вход",
			in:   make(map[string]int),
			out:  []int{},
		},
		{
			name: "Заполненная map[string]int",
			in: map[string]int{
				"1": 1,
				"2": 2,
				"3": 3,
			},
			out: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.out, collections.ValuesOfMap(tt.in))
		})
	}
}

func TestSlice_Join(t *testing.T) {
	tests := []struct {
		name string
		in   [][]int
		out  []int
	}{
		{
			name: "Пустой вход",
			in:   [][]int{},
			out:  []int{},
		},
		{
			name: "Заполненная slice",
			in:   [][]int{{1, 2, 3}, {4, 5, 6}},
			out:  []int{1, 2, 3, 4, 5, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.out, collections.Join(tt.in...))
		})
	}
}

func TestSlice_Single(t *testing.T) {
	tests := []struct {
		name string
		in   int
		out  []int
	}{
		{
			name: "",
			in:   0,
			out:  []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.out, collections.Single(tt.in))
		})
	}
}

func TestSlice_Group(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		fn   func(int) string
		out  map[string][]int
	}{
		{
			name: "Пустой вход",
			in:   []int{},
			fn:   func(i int) string { return "" },
			out:  make(map[string][]int),
		},
		{
			name: "Разбивка на четные и нечетные",
			in:   []int{1, 2, 3, 4, 5, 6, 7, 8},
			fn: func(i int) string {
				if i%2 == 0 {
					return "0"
				}
				return "1"
			},
			out: map[string][]int{
				"0": {2, 4, 6, 8},
				"1": {1, 3, 5, 7},
			},
		},
		{
			name: "Разбивка на остатки 0, 1, 2",
			in:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			fn: func(i int) string {
				if i%3 == 0 {
					return "0"
				} else if i%3 == 1 {
					return "1"
				}
				return "2"
			},
			out: map[string][]int{
				"0": {3, 6, 9},
				"1": {1, 4, 7, 10},
				"2": {2, 5, 8},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := collections.Group(tt.in, tt.fn)
			for key := range res {
				assert.ElementsMatch(t, tt.out[key], res[key])
			}
		})
	}
}

func TestSlice_Dictionary(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		fn   func(int) string
		out  map[string]int
	}{
		{
			name: "Пустой вход",
			in:   []int{},
			fn:   func(i int) string { return "" },
			out:  make(map[string]int),
		},
		{
			name: "Разбивка на остатки 0, 1, 2",
			in:   []int{0, 1, 2},
			fn: func(i int) string {
				if i%3 == 0 {
					return "0"
				} else if i%3 == 1 {
					return "1"
				}
				return "2"
			},
			out: map[string]int{
				"0": 0,
				"1": 1,
				"2": 2,
			},
		},
		{
			name: "Разбивка на остатки 0, 1, 2 с затиранием",
			in:   []int{0, 1, 2, 3, 4, 5},
			fn: func(i int) string {
				if i%3 == 0 {
					return "0"
				} else if i%3 == 1 {
					return "1"
				}
				return "2"
			},
			out: map[string]int{
				"0": 3,
				"1": 4,
				"2": 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := collections.Dictionary(tt.in, tt.fn)
			for key := range res {
				assert.Equal(t, tt.out[key], res[key])
			}
		})
	}
}

func TestSkip(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		size  int
		want  []int
	}{
		{
			name:  "Пропуск половины элементов",
			slice: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			size:  5,
			want:  []int{6, 7, 8, 9, 10},
		},
		{
			name:  "Пропустить количество элементов больше чем в слайсе",
			slice: []int{1, 2, 3, 4, 5},
			size:  10,
			want:  nil,
		},
		{
			name:  "Пропустить 0 элементов",
			slice: []int{1, 2, 3, 4, 5},
			size:  0,
			want:  []int{1, 2, 3, 4, 5},
		},
		{
			name:  "Пропустить элементы в пустом слайсе",
			slice: []int{},
			size:  5,
			want:  nil,
		},
		{
			name:  "Пропустить 0 элементов в пустом слайсе",
			slice: []int{},
			size:  0,
			want:  []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, collections.Skip(tt.slice, tt.size))
		})
	}
}

func TestForEachWithError(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected error
	}{
		{
			name:     "empty slice",
			input:    []int{},
			expected: nil,
		},
		{
			name:     "successful iteration",
			input:    []int{1, 2, 3},
			expected: nil,
		},
		{
			name:     "error during iteration",
			input:    []int{1, 2, 3},
			expected: errors.New("test error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := collections.ForEachWithError(tc.input, func(i int) error {
				if i == 2 {
					return tc.expected
				}
				return nil
			})

			if err != tc.expected {
				t.Errorf("expected error %v, but got %v", tc.expected, err)
			}
		})
	}
}
