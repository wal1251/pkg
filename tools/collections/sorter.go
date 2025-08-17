package collections

import (
	"sort"

	"golang.org/x/exp/constraints"
)

type sorter[T any] struct {
	slice []T
	less  func(T, T) bool
}

func (s *sorter[T]) Len() int {
	return len(s.slice)
}

func (s *sorter[T]) Less(i, j int) bool {
	return s.less(s.slice[i], s.slice[j])
}

func (s *sorter[T]) Swap(i, j int) {
	s.slice[i], s.slice[j] = s.slice[j], s.slice[i]
}

// NewSorter возвращает sort.Interface, для которого порядок сортировки определен функцией less.
func NewSorter[T any](slice []T, less func(T, T) bool) sort.Interface {
	return &sorter[T]{
		slice: slice,
		less:  less,
	}
}

func LessFn[T any, R constraints.Ordered](f func(T) R) func(T, T) bool {
	return func(l T, r T) bool {
		return f(l) < f(r)
	}
}

func SortWith[T any, R constraints.Ordered](s []T, f func(T) R) {
	sort.Sort(NewSorter(s, LessFn(f)))
}
