package collections_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/tools/collections"
)

func TestSyncList(t *testing.T) {
	syncList := collections.NewList[int]()

	var wg sync.WaitGroup

	count := 10
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			syncList.Add(i, i)
		}(i)
	}
	wg.Wait()

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			syncList.Set(i, i)
		}(i)
	}
	wg.Wait()

	for i := 0; i < count; i++ {
		require.Equal(t, i, syncList.Get(i))
	}
}
