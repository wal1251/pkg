package confluent

import (
	"context"

	"github.com/wal1251/pkg/core"
	"github.com/wal1251/pkg/core/bus"
	"github.com/wal1251/pkg/providers/kafka"
	"github.com/wal1251/pkg/tools/collections"
)

var _ bus.EventBus[*kafka.Message] = (*EventBus)(nil)

// EventBus предоставляет реализацию событийной шины bus.EventBus на базе KAFKA.
type EventBus struct {
	prefix      string
	producer    kafka.Producer
	newConsumer func(topic string) (kafka.Consumer, error)
	consumers   []kafka.Consumer
	syncNotify  bool
}

// Notify выполняет публикацию событий в топике name, если топик - пустая строка, тогда топик для публикации будет взят
// из сообщения. Перед публикацией к имени топика добавляется префикс.
func (b *EventBus) Notify(ctx context.Context, name string, events ...*kafka.Message) error {
	return kafka.Send(ctx, b.producer, b.syncNotify, collections.Map(events, func(event *kafka.Message) *kafka.Message {
		msg := *event
		if name == "" {
			msg.Topic = kafka.WithPrefix(b.prefix).Map(msg.Topic)
		} else {
			msg.Topic = kafka.WithPrefix(b.prefix).Map(name)
		}

		return &msg
	})...)
}

// Subscribe регистрирует подписчика топика, к имени топика добавляется префикс. В опубликованных сообщениях топик указан
// без префикса.
func (b *EventBus) Subscribe(ctx context.Context, name string, subscriber bus.Subscriber[*kafka.Message]) error {
	consumer, err := b.newConsumer(kafka.WithPrefix(b.prefix).Map(name))
	if err != nil {
		return err
	}

	return consumer.Subscribe(ctx, bus.SubscriberWith(subscriber, withTopicTransform(kafka.WithoutPrefix(b.prefix))))
}

// Close закрывает соединения, прекращает потребление сообщений.
func (b *EventBus) Close(ctx context.Context) {
	b.producer.Close(ctx)
	collections.ForEach(b.consumers, func(consumer kafka.Consumer) { consumer.Close(ctx) })
}

// NewEventBus возвращает новый экземпляр EventBus. Если в конфиге указан префикс, то данный префикс автоматически
// добавляется к имени топика при публикации.
func NewEventBus(ctx context.Context, cfg *kafka.Config, onError core.ErrorCallback, opts ...kafka.ClientOption) (*EventBus, error) {
	producer, err := NewProducer(cfg, kafka.WithClientID(cfg), kafka.Options(opts...))
	if err != nil {
		return nil, err
	}

	EventsPublish(ctx, producer.Client().Events(), EventsLogger(producer))

	requireACKs, ok := cfg.ProducerConfig[kafka.ClientConfigRequestRequiredACKs]
	isSync := ok && requireACKs != kafka.ClientConfigRequestRequiredACKsNo

	return &EventBus{
		prefix:   cfg.Prefix,
		producer: producer,
		newConsumer: func(topic string) (kafka.Consumer, error) {
			consumer, err := NewConsumer(cfg, collections.Single(topic), onError,
				kafka.WithClientID(cfg), kafka.Options(opts...))
			if err != nil {
				return nil, err
			}

			return consumer, nil
		},
		syncNotify: isSync,
	}, nil
}

func withTopicTransform(transformation core.Map[string, string]) bus.SubscriberMiddleware[*kafka.Message] {
	return bus.SubscriberMiddlewareFn[*kafka.Message](func(ctx context.Context, messages []*kafka.Message, next bus.Subscriber[*kafka.Message]) error {
		collections.ForEach(messages, func(message *kafka.Message) {
			message.Topic = transformation.Map(message.Topic)
		})

		return next.Publish(ctx, messages...)
	})
}
