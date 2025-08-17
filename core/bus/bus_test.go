package bus_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/core/bus"
)

func TestSyncEventBus(t *testing.T) {
	tests := []struct {
		name   string
		events map[string][]string
	}{
		{
			name: "Базовый кейс",
			events: map[string][]string{
				"foo": {"1", "2", "3"},
				"bar": {"4", "5", "6", "7"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()

			results := make(map[string][]string)

			syncBus := bus.NewSyncEventBus[string]()

			for topic := range tt.events {
				currentTopic := topic
				require.NoError(t, syncBus.Subscribe(ctx, topic, bus.SubscriberFn[string](func(_ context.Context, events ...string) error {
					results[currentTopic] = append(results[currentTopic], events...)
					return nil
				})), "must never return error")
			}

			for topic, events := range tt.events {
				require.NoError(t, syncBus.Notify(context.TODO(), topic, events...), "must never return error")
			}

			assert.Equal(t, tt.events, results, "unexpected results")
		})
	}
}
