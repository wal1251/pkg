package bus_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/core/bus"
	"github.com/wal1251/pkg/tools/collections"
)

func TestPublisherWrap(t *testing.T) {
	var results []string

	tests := []struct {
		name        string
		events      []string
		pub         bus.SubscriberFn[string]
		middlewares []bus.SubscriberMiddlewareFn[string]
		want        []string
		wantError   bool
	}{
		{
			name: "Базовый кейс",
			events: []string{
				"foo",
				"bar",
			},
			pub: func(ctx context.Context, messages ...string) error {
				results = append(results, messages...)
				return nil
			},
			middlewares: []bus.SubscriberMiddlewareFn[string]{
				func(ctx context.Context, messages []string, next bus.Subscriber[string]) error {
					results = append(results, "before-1")
					err := next.Publish(ctx, messages...)
					results = append(results, "after-1")
					return err
				},
				func(ctx context.Context, messages []string, next bus.Subscriber[string]) error {
					results = append(results, "before-2")
					err := next.Publish(ctx, messages...)
					results = append(results, "after-2")
					return err
				},
			},
			want: []string{
				"before-1",
				"before-2",
				"foo",
				"bar",
				"after-2",
				"after-1",
			},
		},
		{
			name: "Возврат ошибки в посреднике",
			events: []string{
				"foo",
				"bar",
			},
			pub: func(ctx context.Context, messages ...string) error {
				results = append(results, messages...)
				return nil
			},
			middlewares: []bus.SubscriberMiddlewareFn[string]{
				func(ctx context.Context, messages []string, next bus.Subscriber[string]) error {
					results = append(results, "before-1")
					return errors.New("fake error")
				},
				func(ctx context.Context, messages []string, next bus.Subscriber[string]) error {
					results = append(results, "before-2")
					err := next.Publish(ctx, messages...)
					results = append(results, "after-2")
					return err
				},
			},
			want: []string{
				"before-1",
			},
			wantError: true,
		},
		{
			name: "Нет посредников",
			events: []string{
				"foo",
				"bar",
			},
			pub: func(ctx context.Context, messages ...string) error {
				results = append(results, messages...)
				return nil
			},
			want: []string{
				"foo",
				"bar",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results = make([]string, 0)
			mw := collections.Map(tt.middlewares,
				func(t bus.SubscriberMiddlewareFn[string]) bus.SubscriberMiddleware[string] { return t })

			err := bus.SubscriberWith[string](tt.pub, mw...).Publish(context.TODO(), tt.events...)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, results)
		})
	}
}
