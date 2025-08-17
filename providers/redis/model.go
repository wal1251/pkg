package redis

import (
	"context"

	"github.com/adjust/rmq/v5"
)

var _ QueueConnectionInterface = (*QueueConnectionClient)(nil)

type (
	FailedMessageStrategy func(ctx context.Context, d rmq.Delivery, err error)

	Message struct {
		Value any
	}

	QueueConnectionClient struct {
		rmq rmq.Connection
	}
)
