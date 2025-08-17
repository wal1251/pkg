package concurrent_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/tools/concurrent"
)

func TestDebounce(t *testing.T) {
	var mx sync.Mutex
	result := make([]int, 0)
	expect := make([]int, 0)

	d := concurrent.NewDebounce[int](func(_ context.Context, ts []int) {
		mx.Lock()
		defer mx.Unlock()
		result = append(result, ts...)
	}, 20, 20*time.Millisecond, context.TODO)

	defer func() {
		mx.Lock()
		defer mx.Unlock()
		assert.ElementsMatch(t, expect, result)
	}()
	defer d.Stop()

	for i := 0; i < 100; i++ {
		d.Add(i)
		expect = append(expect, i)
		time.Sleep(2 * time.Millisecond)
	}
}

func TestDebounceAddAfterClose(t *testing.T) {
	d := concurrent.NewDebounce[int](func(_ context.Context, ts []int) { time.Sleep(5 * time.Millisecond) }, 20, 20*time.Millisecond, context.TODO)
	go func() {
		for i := 0; i < 500; i++ {
			d.Add(i)
			time.Sleep(2 * time.Millisecond)
		}
	}()
	time.Sleep(50 * time.Millisecond)
	d.Stop()
}
