package collections

import (
	"fmt"
	"strings"
)

// Set представляет собой коллекцию типа множество - хранит множество уникальных элементов.
type Set[T comparable] map[T]any

// NewSet возвращает новое множество. Добавляет values в качестве элементов множества.
func NewSet[T comparable](values ...T) Set[T] {
	set := make(Set[T])
	set.Add(values...)

	return set
}

// Add добавляет аргументы в множество.
func (s Set[T]) Add(values ...T) {
	for _, value := range values {
		s[value] = struct{}{}
	}
}

// Contains проверяет, содержит ли множество все указанные элементы values, вернет true если добавлен элемент, который
// множество не содержит. Для пустого values всегда возвращает true.
func (s Set[T]) Contains(values ...T) bool {
	for _, value := range values {
		if _, ok := s[value]; !ok {
			return false
		}
	}

	return true
}

func (s Set[T]) NotContains(values ...T) (T, bool) {
	var blank T

	for _, value := range values {
		if _, ok := s[value]; !ok {
			return value, true
		}
	}

	return blank, false
}

// ContainsAny проверяет, содержит ли множество хоть один указанный элемент values, вернет true если добавлен элемент, который
// множество не содержит. Для пустого values всегда возвращает false.
func (s Set[T]) ContainsAny(values ...T) bool {
	for _, value := range values {
		if _, ok := s[value]; ok {
			return true
		}
	}

	return false
}

// Remove удаляет элементы из множества.
func (s Set[T]) Remove(values ...T) {
	for _, value := range values {
		delete(s, value)
	}
}

// Len возвращает количество элементов в множестве.
func (s Set[T]) Len() int {
	return len(map[T]any(s))
}

// ToSlice возвращает слайс элементов, хранимых множеством. Не гарантирует соответствие порядка элементов слайса порядку
// добавления элементов в множество.
func (s Set[T]) ToSlice() []T {
	result := make([]T, 0, len(map[T]any(s)))
	for elem := range s {
		result = append(result, elem)
	}

	return result
}

// ForEach вызывает f для каждого элемента множества. Не гарантирует соответствие порядка обхода элементов
// добавления элементов в множество.
func (s Set[T]) ForEach(act func(T)) {
	for elem := range s {
		act(elem)
	}
}

// String возвращает строку с элементами множества через запятую.
func (s Set[T]) String() string {
	return fmt.Sprintf("[%s]",
		strings.Join(Map(s.ToSlice(), func(t T) string { return fmt.Sprint(t) }), ", "))
}

// CountUnique возвращает количество уникальных элементов в collection.
func CountUnique[T comparable](collection []T) int {
	return NewSet[T](collection...).Len()
}

// Contains возвращает true если collection содержит значение v.
func Contains[T comparable](collection []T, values ...T) bool {
	return NewSet(collection...).Contains(values...)
}
