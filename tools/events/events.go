package events

import (
	"context"

	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/providers/kafka"
)

// EventRead This function is used to read events from a Kafka message and convert the message value to a typed object
// of any type T. The function takes an accept function as a parameter, which is used to process the typed object read
// from the message. If there is an error reading the message or converting the message value to a typed object,
// an error message is logged. The function returns a function that takes a context and a Kafka  message as parameters.
// This function is used to read the message and call the accept function to process the typed object.
func EventRead[T any](accept func(context.Context, T)) func(context.Context, *kafka.Message) {
	return func(ctx context.Context, msg *kafka.Message) {
		value, err := kafka.ValueJSONToTyped[T](msg.Value)
		if err != nil {
			logs.FromContext(ctx).Err(err).Msgf("failed to read message")

			return
		}

		accept(ctx, value)
	}
}
