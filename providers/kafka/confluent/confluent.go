// Package confluent содержит реализацию клиентов KAFKA, которая базируется на библиотеке confluent-kafka-go.
package confluent

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"

	confluent "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/wal1251/pkg/core"
	"github.com/wal1251/pkg/core/bus"
	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/providers/kafka"
	"github.com/wal1251/pkg/tools/collections"
)

var _ bus.Subscriber[confluent.Event] = (*StatsHolder)(nil)

// StatsHolder для хранения статистики клиента, реализует интерфейс bus.Subscriber для публикации данных статистики.
type StatsHolder atomic.Value

// Publish см. bus.Subscriber.
func (h *StatsHolder) Publish(_ context.Context, events ...confluent.Event) error {
	for _, event := range events {
		if stats, ok := event.(*confluent.Stats); ok {
			(*atomic.Value)(h).Store([]byte(stats.String()))
		}
	}

	return nil
}

// Demand см. bus.Subscriber.
func (h *StatsHolder) Demand() int {
	return math.MaxInt
}

func (h *StatsHolder) String() string {
	return string(h.Get())
}

// Get возвращает текущее состояние статистики в формате JSON.
func (h *StatsHolder) Get() []byte {
	if stats := (*atomic.Value)(h).Load(); stats != nil {
		if raw, ok := stats.([]byte); ok {
			return raw
		}
	}

	return nil
}

// ReadHeaders возвращает представление заголовка сообщения KAFKA.
func ReadHeaders(message *confluent.Message) collections.MultiMap[string, []byte] {
	headers := make(collections.MultiMap[string, []byte])
	for _, header := range message.Headers {
		headers.Append(header.Key, header.Value)
	}

	return headers
}

// MakeHeaders возвращает новый предварительно заполненный заголовок сообщения KAFKA.
func MakeHeaders(source collections.MultiMap[string, []byte]) []confluent.Header {
	var headers []confluent.Header
	for key, value := range source {
		if headers == nil {
			headers = make([]confluent.Header, 0, len(headers))
		}

		for _, v := range value {
			headers = append(headers, confluent.Header{Key: key, Value: v})
		}
	}

	return headers
}

// MakeMessage создает новое сообщение KAFKA из переданного представления. Если сообщение не удалось создать, вернет
// ошибку.
func MakeMessage(message *kafka.Message) (*confluent.Message, error) {
	partition := confluent.PartitionAny
	if message.Partition != nil {
		partition = *message.Partition
	}

	value, err := message.Value.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare message value: %w", err)
	}

	return &confluent.Message{
		TopicPartition: confluent.TopicPartition{
			Topic:     &message.Topic,
			Partition: partition,
		},
		Key:     message.Key,
		Value:   value,
		Headers: MakeHeaders(message.Headers),
	}, nil
}

// ReadMessage возвращает представление сообщения KAFKA.
func ReadMessage(message *confluent.Message) *kafka.Message {
	newMessage := &kafka.Message{
		Partition: &message.TopicPartition.Partition,
		Headers:   ReadHeaders(message),
	}

	return newMessage.
		WithValue(kafka.ValueBytes(message.Value)).
		WithKeyBytes(message.Key)
}

// DefaultMessageReader возвращает функцию преобразования bus.Map сообщения KAFKA к каноничному виду kafka.Message для пакета.
func DefaultMessageReader(consumer *confluent.Consumer, callback core.ErrorCallback) core.Map[*confluent.Message, *kafka.Message] {
	return func(msg *confluent.Message) *kafka.Message {
		return ReadMessage(msg).WithAck(func(ctx context.Context) {
			if err := CommitMessage(consumer, msg); err != nil {
				logs.FromContext(ctx).Err(err).Msg("commit failed")
				callback.OnError(err)
			}
		})
	}
}

// EventsPublish публикует указанный канал сообщений клиента KAFKA на заданных подписчиков.
func EventsPublish(ctx context.Context, events <-chan confluent.Event, subscribers ...bus.Subscriber[confluent.Event]) {
	logger := logs.FromContext(ctx)
	pub := bus.SubscribeAll(subscribers...)

	done := ctx.Done()
	go func() {
		for event := range events {
			select {
			case <-done:
				return
			default:
				if err := pub.Publish(ctx, event); err != nil {
					logger.Err(err).Msg("failed to publish kafka event")
				}
			}
		}
	}()
}

// EventsLogger возвращает подписчик на события KAFKA, логирует публикуемые в него события KAFKA.
func EventsLogger(component core.Component) bus.SubscriberFn[confluent.Event] {
	return func(ctx context.Context, events ...confluent.Event) error {
		log := logs.LocalContext(component).To(logs.FromContext(ctx).Debug)
		for _, event := range events {
			if _, ok := event.(*confluent.Stats); !ok {
				log.Msgf("kafka event: %v", event)
			}
		}

		return nil
	}
}

// NewStatsHolder возвращает новый StatsHolder, который так же является подписчиком на события KAFKA, при публикации
// событий статистики на него, обновляет свое состояние. Статистику можно считать методом StatsHolder.Get().
func NewStatsHolder() *StatsHolder {
	return new(StatsHolder)
}

// ConfigMap создает конфигурацию клиента KAFKA.
func ConfigMap(config kafka.ClientConfig) *confluent.ConfigMap {
	configMap := make(confluent.ConfigMap)
	for key, value := range config {
		configMap[key] = value
	}

	return &configMap
}

// CommitMessage фиксирует смещения по указанным сообщениям с помощью клиента confluent.Consumer.
func CommitMessage(consumer *confluent.Consumer, messages ...*confluent.Message) error {
	for _, msg := range messages {
		offsets, err := consumer.CommitMessage(msg)
		if err != nil {
			return fmt.Errorf("%v can't commit message %v: %w", consumer, messages, err)
		}

		for _, offset := range offsets {
			if offset.Error != nil {
				return fmt.Errorf("%v can't commit %v: %w", consumer, offset, err)
			}
		}
	}

	return nil
}

// MakeTopicSpec возвращает новую спецификацию топика KAFKA.
func MakeTopicSpec(topic kafka.TopicMetadata) confluent.TopicSpecification {
	return confluent.TopicSpecification{
		Topic:             topic.Name,
		NumPartitions:     topic.Partitions,
		ReplicationFactor: topic.ReplicationFactor,
	}
}
