package singleton

import "sync"

// Singleton представляет структуру Одиночки, хранящую единственный экземпляр объекта типа T.
type Singleton[T any] struct {
	instance *T
	once     sync.Once
	initFunc func() *T
}

// NewSingleton создает и возвращает новый экземпляр Singleton с предоставленной функцией инициализации.
// Эта функция будет вызвана не более одного раза для создания экземпляра.
//
// Параметры:
// initFunc - функция, возвращающая новый экземпляр объекта типа *T.
//
// Возвращаемое значение:
// Возвращает новый экземпляр Singleton.
func NewSingleton[T any](initFunc func() *T) *Singleton[T] {
	return &Singleton[T]{
		initFunc: initFunc,
	}
}

// Get возвращает единственный экземпляр объекта типа T.
// Если экземпляр еще не создан, он будет создан с использованием функции инициализации, предоставленной в NewSingleton.
// Последующие вызовы Get будут возвращать тот же экземпляр.
//
// Возвращаемое значение:
// Возвращает указатель на экземпляр объекта типа T.
func (s *Singleton[T]) Get() *T {
	s.once.Do(func() {
		s.instance = s.initFunc()
	})

	return s.instance
}
