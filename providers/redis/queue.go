package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/adjust/rmq/v5"
	rv9 "github.com/redis/go-redis/v9"

	"github.com/wal1251/pkg/core/bus"
	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/tools/collections"
	"github.com/wal1251/pkg/tools/serial"
)

type (
	QueueConnectionInterface interface {
		// QueueStats функция позволяет получить rmq.Stats(статистику) для всех открытых очередей.
		QueueStats(_ context.Context) (rmq.Stats, error)

		// GetOpenQueues функция позволяет получить список названий открытых очередей.
		GetOpenQueues() ([]string, error)

		// OpenQueue функция позволяет получить очередь по имени.
		OpenQueue(name string) (rmq.Queue, error)

		// StopAllConsuming функция позволяет остановить использование всех очередей, открытых в этом соединении.
		StopAllConsuming()
	}
)

// NewQueueConnection создает клиент для работы с очередями.
func NewQueueConnection(ctx context.Context, client *rv9.Client, tag string) (QueueConnectionClient, error) {
	logger := logs.FromContext(ctx)
	backgroundErr := make(chan error)

	conn, err := rmq.OpenConnectionWithRedisClient(tag, client, backgroundErr)
	if err != nil {
		return QueueConnectionClient{}, fmt.Errorf("can't open redis connection: %w", err)
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

	return QueueConnectionClient{rmq: conn}, nil
}

func (client *QueueConnectionClient) QueueStats(_ context.Context) (rmq.Stats, error) {
	list, err := client.GetOpenQueues()
	if err != nil {
		return rmq.Stats{}, fmt.Errorf("unable to get opened queues: %w", err)
	}

	stats, err := client.rmq.CollectStats(list)
	if err != nil {
		return rmq.Stats{}, fmt.Errorf("unable to get queues stats: %w", err)
	}

	return stats, nil
}

func (client *QueueConnectionClient) GetOpenQueues() ([]string, error) {
	queues, err := client.rmq.GetOpenQueues()
	if err != nil {
		return nil, fmt.Errorf("can't get opened queues: %w", err)
	}

	return queues, nil
}

func (client *QueueConnectionClient) OpenQueue(name string) (rmq.Queue, error) {
	queue, err := client.rmq.OpenQueue(name)
	if err != nil {
		return nil, fmt.Errorf("can't open queue %s: %w", name, err)
	}

	return queue, nil
}

func (client *QueueConnectionClient) StopAllConsuming() {
	client.rmq.StopAllConsuming()
}

// MakeSubscriber создает подписчика очереди, для публикации событий в эту очередь.
func MakeSubscriber[T any](queue rmq.Queue) bus.SubscriberFn[T] {
	return func(ctx context.Context, messages ...T) error {
		logger := logs.FromContext(ctx)

		data, err := collections.MapWithErr(messages, func(msg T) ([]byte, error) {
			return serial.ToBytes(msg, serial.JSONEncode[T])
		})
		if err != nil {
			logger.Err(err).Msg("unable to marshal message for publishing")

			return err
		}

		if err = queue.PublishBytes(data...); err != nil {
			logger.Err(err).Msg("unable to publish messages")

			return fmt.Errorf("can't publish message to queue: %w", err)
		}

		return nil
	}
}

// MakePublisher вернет издателя очереди, с помощью которого можно подписаться на события очереди.
func MakePublisher[T any](queue rmq.Queue, tag string, onFail FailedMessageStrategy) bus.PublisherFn[T] {
	if onFail == nil {
		onFail = FailedMessageIgnoreStrategy()
	}

	return func(ctx context.Context, subscriber bus.Subscriber[T]) error {
		if _, err := queue.AddConsumerFunc(tag, func(delivery rmq.Delivery) {
			var message T
			var err error

			defer func() {
				if err != nil {
					onFail(ctx, delivery, err)

					return
				}

				if err = delivery.Ack(); err != nil {
					logs.FromContext(ctx).Err(err).Msg("unable to ack message")
				}
			}()

			if message, err = serial.FromBytes([]byte(delivery.Payload()), serial.JSONDecode[T]); err == nil {
				err = subscriber.Publish(ctx, message)
			}
		}); err != nil {
			return fmt.Errorf("can't consume queue message: %w", err)
		}

		return nil
	}
}

func FailedMessageIgnoreStrategy() FailedMessageStrategy {
	return func(ctx context.Context, d rmq.Delivery, err error) {
		logs.FromContext(ctx).Err(err).Msg("failed to process queue message")

		if err = d.Ack(); err != nil {
			logs.FromContext(ctx).Err(err).Msg("unable to ack message")
		}
	}
}

func FailedMessageRejectStrategy() FailedMessageStrategy {
	return func(ctx context.Context, d rmq.Delivery, err error) {
		logs.FromContext(ctx).Err(err).Msg("failed to process queue message")

		if err = d.Reject(); err != nil {
			logs.FromContext(ctx).Err(err).Msg("unable to reject message")
		}
	}
}
