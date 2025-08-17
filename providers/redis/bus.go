package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/adjust/rmq/v5"
	rv9 "github.com/redis/go-redis/v9"

	"github.com/wal1251/pkg/core/bus"
	"github.com/wal1251/pkg/core/logs"
)

var _ bus.EventBus[*Message] = (*EventBus)(nil)

type EventBus struct {
	tag      string
	strategy FailedMessageStrategy
	client   rmq.Connection
	config   *BusConfig
}

func (e *EventBus) Notify(ctx context.Context, topic string, events ...*Message) error {
	queue, err := e.client.OpenQueue(topic)
	if err != nil {
		return fmt.Errorf("can't open queue %s: %w", topic, err)
	}

	return MakeSubscriber[*Message](queue).Publish(ctx, events...)
}

func (e *EventBus) Subscribe(ctx context.Context, topic string, subscriber bus.Subscriber[*Message]) error {
	queue, err := e.client.OpenQueue(topic)
	if err != nil {
		return fmt.Errorf("can't open queue %s: %w", topic, err)
	}
	err = queue.StartConsuming(int64(e.config.ConsumerPrefetchLimit), e.config.ConsumerPollDuration)
	if err != nil {
		return fmt.Errorf("can't start consuming %s: %w", topic, err)
	}

	return MakePublisher[*Message](queue, e.tag, e.strategy).Subscribe(ctx, subscriber)
}

func (e *EventBus) Close(_ context.Context) {
	e.client.StopAllConsuming()
}

func NewEventBus(
	ctx context.Context,
	cfg *BusConfig,
	client *rv9.Client,
	tag string,
	strategy FailedMessageStrategy,
) (*EventBus, error) {
	logger := logs.FromContext(ctx)
	backgroundErr := make(chan error)

	conn, err := rmq.OpenConnectionWithRedisClient(tag, client, backgroundErr)
	if err != nil {
		return nil, fmt.Errorf("can't open redis connection: %w", err)
	}

	cleaner := rmq.NewCleaner(conn)

	go func() {
		for {
			logger.Err(<-backgroundErr).Msg("have fail in redis queue background logic")
		}
	}()
	go func() {
		for range time.Tick(QueueConnectionCleanup) {
			returned, err := cleaner.Clean()
			if err != nil {
				logger.Err(err).Msg("failed to clean up queue connection")

				continue
			}

			if returned != 0 {
				logger.Debug().Msgf("cleaned: %d", returned)
			}
		}
	}()

	return &EventBus{
		client:   conn,
		tag:      tag,
		strategy: strategy,
		config:   cfg,
	}, nil
}
