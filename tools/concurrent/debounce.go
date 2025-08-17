package concurrent

import (
	"context"
	"sync"
	"time"

	"github.com/wal1251/pkg/tools/collections"
)

const debounceMultiplexer = 4

type Debounce[T any] struct {
	batchSize  int
	timeout    time.Duration
	wg         sync.WaitGroup
	items      chan T
	consumer   func([]T)
	closedLock sync.Mutex
	closedFlag bool
}

func (d *Debounce[T]) Add(items ...T) bool {
	d.closedLock.Lock()
	defer d.closedLock.Unlock()

	active := !d.closedFlag
	if active {
		collections.ForEach(items, func(t T) { d.items <- t })
	}

	return active
}

func (d *Debounce[T]) Stop() {
	defer d.wg.Wait()
	func() {
		d.closedLock.Lock()
		defer d.closedLock.Unlock()

		if !d.closedFlag {
			d.closedFlag = true
			close(d.items)
		}
	}()
}

func (d *Debounce[T]) run() {
	tick := time.NewTicker(d.timeout)

	d.wg.Add(1)
	go func() {
		var items []T

		defer d.wg.Done()

		send := false
		for active := true; active; {
			select {
			case item, ok := <-d.items:
				active = ok
				if ok {
					items = append(items, item)
				}

				send = !ok || len(items) >= d.batchSize
			case <-tick.C:
				send = true
			}

			if send {
				send = false

				if len(items) != 0 {
					d.wg.Add(1)

					go func(s []T) {
						defer d.wg.Done()
						d.consumer(s)
					}(append([]T{}, items...))

					items = items[:0]
				}
			}
		}
	}()
}

func NewDebounce[T any](consumer func(context.Context, []T), size int, interval time.Duration, ctx func() context.Context) *Debounce[T] {
	debounce := &Debounce[T]{
		items: make(chan T, debounceMultiplexer*size),
		consumer: func(ts []T) {
			if ctx == nil {
				consumer(context.Background(), ts)
			} else {
				consumer(ctx(), ts)
			}
		},
		batchSize: size,
		timeout:   interval,
	}
	debounce.run()

	return debounce
}
