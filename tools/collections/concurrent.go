package collections

import "sync"

// SyncList потокобезопасный список с основными операциями.
type SyncList[T any] struct {
	sync.RWMutex
	list []T
}

// NewList конструктор списка по входным значениям.
func NewList[T any](list ...T) *SyncList[T] {
	return &SyncList[T]{
		list: list,
	}
}

// Copy получить копию.
func (l *SyncList[T]) Copy() []T {
	l.RLock()
	defer l.RUnlock()

	dest := make([]T, len(l.list))
	copy(dest, l.list)

	return dest
}

// Get получить элемент по индексу.
func (l *SyncList[T]) Get(i int) T {
	l.RLock()
	defer l.RUnlock()

	return l.list[i]
}

// Set Записать в ячейку с индексом i значение.
func (l *SyncList[T]) Set(i int, v T) {
	l.Lock()
	defer l.Unlock()

	l.list[i] = v
}

// Add добавить элемент в конец списка.
func (l *SyncList[T]) Add(v ...T) {
	l.Lock()
	defer l.Unlock()

	l.list = append(l.list, v...)
}
