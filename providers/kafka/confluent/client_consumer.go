package confluent

import (
	"context"
	"fmt"
	"sync/atomic"

	confluent "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/wal1251/pkg/core"
	"github.com/wal1251/pkg/core/bus"
	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/providers/kafka"
)

const (
	ConsumerComponentName = "confluent-kafka-go-consumer" // Имя компонента клиента-потребителя.

	ClientConfigReBalanceEnable = "go.application.rebalance.enable" // Включение обработки сообщений балансировки.
)

const (
	ConsumerStateIdle     = iota // Состояние потребителя: Покой.
	ConsumerStateActive          // Состояние потребителя: Активен.
	ConsumerStateDisposed        // Состояние потребителя: Остановлен.
)

var _ kafka.Consumer = (*Consumer)(nil)

type (
	// Consumer реализация клиента-потребителя KAFKA, реализует kafka.Consumer. Клиент осуществляет только цикл опроса
	// брокера. Стратегия подтверждения сообщений определяется на уровне компонента kafka.Poller.
	Consumer struct {
		core.Component

		isClosed  chan any
		topics    []string
		state     int32
		consumer  *confluent.Consumer
		newPoller func(consumer *confluent.Consumer, subscriber bus.Subscriber[*kafka.Message]) kafka.Poller
	}
)

// Subscribe запускает цикл опроса брокера, полученные сообщения будут публиковаться на указанного потребителя. Если
// переданный при подписке контекст будет отменен - цикл опроса будет завершен.
func (c *Consumer) Subscribe(ctx context.Context, subscriber bus.Subscriber[*kafka.Message]) error {
	logger := logs.FromContext(ctx)
	log := logs.LocalContext(c).To

	if !atomic.CompareAndSwapInt32(&c.state, ConsumerStateIdle, ConsumerStateActive) {
		state := atomic.LoadInt32(&c.state)

		return fmt.Errorf("%w: state: %d", kafka.ErrInvalidConsumerState, state)
	}

	if err := c.consumer.SubscribeTopics(c.topics, nil); err != nil {
		return fmt.Errorf("failed to subscribe kafka consumer %v: %w", c.topics, err)
	}

	log(logger.Info).Msgf("subscribed: %v", c.topics)

	poller := c.newPoller(c.consumer, subscriber)

	go func(done <-chan struct{}) {
		defer c.dispose(ctx)

		log(logger.Info).Msg("start consuming")

		for c.IsActive() {
			select {
			case <-done:
				log(logger.Info).Msg("stopping to consume due to context cancel")

				return

			default:
				if err := poller.Poll(ctx); err != nil {
					log(logger.Error).Err(err).Msg("stopping to consume due to consumer fatal error")

					return
				}
			}
		}
	}(ctx.Done())

	return nil
}

// Close останавливает цикл потребления и освобождает занятые клиентом ресурсы.
func (c *Consumer) Close(ctx context.Context) {
	c.dispose(ctx)
	<-c.isClosed
}

func (c *Consumer) String() string {
	return c.consumer.String()
}

// IsActive вернет true, если цикл опроса клиента активен.
func (c *Consumer) IsActive() bool {
	return atomic.LoadInt32(&c.state) == ConsumerStateActive
}

// Client вернет клиент-потребитель более низкого уровня.
func (c *Consumer) Client() *confluent.Consumer {
	return c.consumer
}

// ConfigurePoller устанавливает конструктор компонента опроса брокера. Позволяет заменить конструктор по умолчанию на
// кастомную реализацию. Позволит реализовать особые стратегии опроса и подтверждения сообщений. Должен быть вызван
// перед вызовом метода Subscribe().
// Например, если мы установим опцию авто-фиксации смещений на уровне клиента, то можно передать конструктор компоненты,
// которая осуществляет фиксацию смещений.
func (c *Consumer) ConfigurePoller(newPoller func(consumer *confluent.Consumer, subscriber bus.Subscriber[*kafka.Message]) kafka.Poller) {
	c.newPoller = newPoller
}

func (c *Consumer) dispose(ctx context.Context) {
	logger := logs.FromContext(ctx)
	log := logs.LocalContext(c).To

	if !atomic.CompareAndSwapInt32(&c.state, ConsumerStateIdle, ConsumerStateDisposed) &&
		!atomic.CompareAndSwapInt32(&c.state, ConsumerStateActive, ConsumerStateDisposed) {
		return
	}

	if err := c.consumer.Unsubscribe(); err != nil {
		log(logger.Error).Err(err).Msg("failed to unsubscribe consumer")
	}

	if err := c.consumer.Close(); err != nil {
		log(logger.Error).Err(err).Msg("failed to close consumer")
	}

	close(c.isClosed)

	log(logger.Info).Msg("consumer closed")
}

// NewConsumer возвращает новый клиент-потребитель Consumer.
func NewConsumer(cfg *kafka.Config, topics []string, onError core.ErrorCallback, opts ...kafka.ClientOption) (*Consumer, error) {
	clientConfig, err := kafka.ConsumerConfig(cfg, opts...)
	if err != nil {
		return nil, fmt.Errorf("can't create kafka consumer: %w", err)
	}

	consumer, err := confluent.NewConsumer(ConfigMap(clientConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	return &Consumer{
		Component: core.NewDefaultComponent(ConsumerComponentName, consumer.String()),
		topics:    topics,
		consumer:  consumer,
		newPoller: PollerDefaultFactory(cfg, onError),
		isClosed:  make(chan any),
	}, nil
}
