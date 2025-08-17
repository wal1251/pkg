package concurrent

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultBufferMultiply = 2
)

var errTimeOut = errors.New("timeout error")

type WorkersPool struct {
	lock         sync.Mutex
	cond         *sync.Cond
	wg           sync.WaitGroup
	once         sync.Once
	jobs         []func() error
	workersCount int
	bufferSize   int
	queueCount   int32
	err          error
	cancel       func()
	finished     chan struct{}
	isClosed     bool
}

func (p *WorkersPool) PendingCount() int {
	return int(atomic.LoadInt32(&p.queueCount))
}

func (p *WorkersPool) BufferedCount() int {
	p.lock.Lock()
	defer p.lock.Unlock()

	return len(p.jobs)
}

func (p *WorkersPool) IsClosed() bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.isClosed
}

func (p *WorkersPool) Add(job func() error) {
	atomic.AddInt32(&p.queueCount, 1)
	p.wg.Add(1)
	go func() {
		p.cond.L.Lock()
		defer p.cond.L.Unlock()
		for len(p.jobs) >= p.bufferSize && !p.isClosed {
			p.cond.Wait()
		}
		atomic.AddInt32(&p.queueCount, -1)
		if p.isClosed {
			p.wg.Done()

			return
		}
		p.jobs = append(p.jobs, func() error {
			defer p.wg.Done()

			return job()
		})
		p.cond.Broadcast()
	}()
}

func (p *WorkersPool) run(ctx context.Context) {
	next := func() func() error {
		p.cond.L.Lock()
		defer p.cond.L.Unlock()
		for len(p.jobs) == 0 && !p.isClosed {
			p.cond.Wait()
		}
		if p.isClosed {
			return nil
		}
		job := p.jobs[0]
		p.jobs = p.jobs[1:]

		p.cond.Broadcast()

		return job
	}

	for i := 0; i < p.workersCount; i++ {
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					job := next()
					if job == nil {
						break
					}
					if err := job(); err != nil && p.cancel != nil {
						p.once.Do(func() {
							p.err = err
							p.finished <- struct{}{}
							p.cancel()
						})
					}
				}
			}
		}(ctx)
	}
}

// RunIgnoreError игнорирует ошибки, если какой-то воркер завершился с ошибкой.
func (p *WorkersPool) RunIgnoreError(ctx context.Context) {
	p.run(ctx)
}

// RunStrictError если какой-то воркер завершается с ошибкой, то завершить работу worker poll, вернет первую не nil ошибку.
func (p *WorkersPool) RunStrictError(ctx context.Context) {
	childCtx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	p.run(childCtx)
}

func (p *WorkersPool) Wait(timeout time.Duration) error {
	go func() {
		p.wg.Wait()
		p.finished <- struct{}{}
	}()

	select {
	case <-p.finished:
		return p.err
	case <-time.After(timeout):
		p.cancel()

		return errTimeOut
	}
}

func (p *WorkersPool) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.isClosed = true
}

func NewWorkersPool(workersCount int) *WorkersPool {
	workerPool := &WorkersPool{
		jobs:         make([]func() error, 0, defaultBufferMultiply*workersCount),
		finished:     make(chan struct{}),
		workersCount: workersCount,
		bufferSize:   defaultBufferMultiply * workersCount,
	}
	workerPool.cond = sync.NewCond(&workerPool.lock)

	return workerPool
}
