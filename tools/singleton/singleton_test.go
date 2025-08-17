package singleton

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSingleton_Get_Concurrency проверяет корректность работы Singleton при многопоточном доступе.
func TestSingleton_Get_Concurrency(t *testing.T) {
	var initCount int
	initFunc := func() *int {
		initCount++
		return &initCount
	}

	singleton := NewSingleton(initFunc)

	var wg sync.WaitGroup
	const goroutines = 30 // Количество горутин для теста

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			instance := singleton.Get()
			assert.Equal(t, 1, *instance)
		}()
	}

	wg.Wait()
	// Проверяем, что инициализация произошла только один раз
	assert.Equal(t, 1, initCount)
}

// TestSingleton_Get_DifferentTypes проверяет Singleton для различных типов.
func TestSingleton_Get_DifferentTypes(t *testing.T) {
	initFuncInt := func() *int {
		num := 5
		return &num
	}
	initFuncString := func() *string {
		str := "singleton"
		return &str
	}

	singletonInt := NewSingleton(initFuncInt)
	singletonString := NewSingleton(initFuncString)

	instanceInt := singletonInt.Get()
	instanceString := singletonString.Get()

	assert.Equal(t, 5, *instanceInt)
	assert.Equal(t, "singleton", *instanceString)
}

// TestSingleton_Get_NoInitFunc проверяет поведение Singleton без функции инициализации.
func TestSingleton_Get_NoInitFunc(t *testing.T) {
	var singleton *Singleton[int] = NewSingleton[int](nil)

	assert.Panics(t, func() {
		_ = singleton.Get()
	}, "Singleton should panic if initFunc is nil")
}
