package events

import (
	"context"

	"github.com/google/uuid"

	"github.com/wal1251/pkg/core/bus"
	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/core/presenters"
)

func Logger(view presenters.ViewType, options presenters.ViewOptions) bus.SubscriberFn[Event] {
	return func(ctx context.Context, events ...Event) error {
		logger := logs.FromContext(ctx)

		for _, event := range events {
			attributes := presenters.ParameterView(event.Attributes, view, options)
			if event.ID == uuid.Nil {
				logger.Trace().Msgf("%s: %s %s", event.Type, event.Op, attributes)
			} else {
				logger.Trace().Msgf("%s(%v): %s %s", event.Type, event.ID, event.Op, attributes)
			}
		}

		return nil
	}
}
