package concurrent_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/tools/concurrent"
)

func TestWorkerPool(t *testing.T) {
	type Task struct {
		Id          uuid.UUID
		IsCompleted bool
	}

	listTask := []Task{
		{Id: uuid.New()},
		{Id: uuid.New()},
		{Id: uuid.New()},
	}

	workerPool := concurrent.NewWorkersPool(5)

	for index := range listTask {
		func(i int) {
			workerPool.Add(func() error {
				listTask[i].IsCompleted = true
				return nil
			})
		}(index)
	}

	checker := func() bool {
		result := true
		for _, task := range listTask {
			result = result && task.IsCompleted
		}
		return result
	}

	workerPool.RunIgnoreError(context.TODO())

	err := workerPool.Wait(10 * time.Second)
	require.NoError(t, err)

	require.True(t, checker())
}

func TestWorkersPool_SuccessWithErrors(t *testing.T) {
	type Task struct {
		Id          int
		IsCompleted bool
	}

	listTask := []Task{
		{Id: 0},
		{Id: 1},
	}

	workerPool := concurrent.NewWorkersPool(5)

	for index := range listTask {
		func(i int) {
			workerPool.Add(func() error {
				if listTask[i].Id%2 != 0 {
					return errors.New("error")
				}
				listTask[i].IsCompleted = true
				return nil
			})
		}(index)
	}

	workerPool.RunIgnoreError(context.TODO())

	err := workerPool.Wait(time.Minute)
	require.NoError(t, err)

	require.True(t, listTask[0].IsCompleted)
	require.False(t, listTask[1].IsCompleted)
}

func TestWorkersPool_Rejected(t *testing.T) {
	urls := []string{
		"https://www.success.com",
		"https://www.rejected.com",
		"https://www.success.com",
		"https://www.success.com",
		"https://www.success.com",
		"https://www.success.com",
	}

	ping := func(url string) error {
		if strings.Contains(url, "rejected") {
			return errors.New("error")
		}
		return nil
	}

	workerPool := concurrent.NewWorkersPool(5)

	for _, elem := range urls {
		func(url string) {
			workerPool.Add(func() error {
				err := ping(url)
				if err != nil {
					return err
				}
				return nil
			})
		}(elem)
	}

	workerPool.RunStrictError(context.TODO())

	err := workerPool.Wait(time.Minute)
	require.Error(t, err)
}

func TestWorkersPool_Timeout(t *testing.T) {
	workerPool := concurrent.NewWorkersPool(5)

	workerPool.Add(func() error {
		time.Sleep(time.Minute)
		return nil
	})

	workerPool.RunStrictError(context.TODO())

	err := workerPool.Wait(10 * time.Second)
	require.Error(t, err)
}
