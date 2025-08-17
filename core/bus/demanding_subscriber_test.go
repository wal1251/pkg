package bus_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/core/bus"
)

func TestAsyncDemandingSubscriber(t *testing.T) {
	var wg sync.WaitGroup
	var results []string

	tests := []struct {
		name       string
		prefetch   int
		messages   []string
		subscriber bus.SubscriberFn[string]
		want       []string
		wantError  bool
	}{
		{
			name:     "Базовый кейс",
			prefetch: 5,
			messages: []string{"1", "2", "3", "4", "5"},
			subscriber: func(ctx context.Context, events ...string) error {
				defer wg.Done()

				if events[0] == "1" { // Остановимся, чтобы померить demand
					time.Sleep(50 * time.Millisecond)
				}

				results = append(results, events...)

				return nil
			},
			want: []string{"1", "2", "3", "4", "5"},
		},
		{
			name:     "Публикуем больше чем емкость",
			prefetch: 5,
			messages: []string{"1", "2", "3", "4", "5", "6", "7"},
			subscriber: func(ctx context.Context, events ...string) error {
				defer wg.Done()

				if events[0] == "1" { // Остановимся, чтобы померить demand
					time.Sleep(50 * time.Millisecond)
				}

				results = append(results, events...)

				return nil
			},
			want: []string{"1", "2", "3", "4", "5", "6", "7"},
		},
		{
			name:     "Вернем ошибку",
			prefetch: 5,
			messages: []string{"1", "2", "3", "4", "5"},
			subscriber: func(ctx context.Context, events ...string) error {
				defer wg.Done()

				if events[0] == "1" { // Остановимся, чтобы померить demand
					time.Sleep(50 * time.Millisecond)
					return errors.New("fake error")
				}

				results = append(results, events...)

				return nil
			},
			want:      []string{"2", "3", "4", "5"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var caughtErr error

			results = make([]string, 0)
			ctx := context.TODO()

			sub := bus.NewAsyncDemandingSubscriber[string](ctx, tt.subscriber, tt.prefetch,
				func(err error) {
					require.Error(t, err, "no error caught with onError callback")
					caughtErr = err
				})

			assert.Equal(t, tt.prefetch, sub.Demand(),
				"initial demand differs from expected prefetch value")

			wg.Add(len(tt.messages))
			require.NoError(t, sub.Publish(ctx, tt.messages...),
				"async subscriber must never return error")

			expectedDemand := tt.prefetch - len(tt.messages)
			if expectedDemand < 0 {
				expectedDemand = 0
			}

			assert.Equal(t, expectedDemand, sub.Demand(),
				"demand differs from expected, check subscriber in argument is waiting")

			wg.Wait()

			assert.Equal(t, tt.prefetch, sub.Demand(), "demand must return to initial capacity")
			assert.Equal(t, tt.want, results,
				"result is differs from expected, check check subscriber in argument is writing result")

			if tt.wantError {
				assert.Error(t, caughtErr,
					"error is expected but not caught, check subscriber returns error")
			} else {
				assert.NoError(t, caughtErr,
					"error is not expected but error callback is called")
			}
		})
	}
}
