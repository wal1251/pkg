package doubleconfluent

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"

	confluent "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/wal1251/pkg/core/bus"
	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/providers/kafka"
	"github.com/wal1251/pkg/tools/collections"
)

// MessageEvent обертка для kafka.Message, реализующая интерфейс confluent.Event.
type MessageEvent struct {
	*kafka.Message
}

// String реализует метод интерфейса confluent.Event.
func (m *MessageEvent) String() string {
	if m.Value != nil { //nolint:typecheck
		return string(m.Value.Must()) //nolint:typecheck
	}

	return ""
}

// StatsHolderMock для хранения статистики клиента, реализует интерфейс bus.Subscriber для публикации данных статистики.
type StatsHolderMock struct {
	stats atomic.Value
}

// NewStatsHolderMock создает новый мок StatsHolder.
func NewStatsHolderMock() *StatsHolderMock {
	return &StatsHolderMock{}
}

func (h *StatsHolderMock) Publish(_ context.Context, events ...confluent.Event) error {
	for _, event := range events {
		if stats, ok := event.(*confluent.Stats); ok {
			h.stats.Store([]byte(stats.String()))
		}
	}

	return nil
}

func (h *StatsHolderMock) Demand() int {
	return math.MaxInt
}

func (h *StatsHolderMock) String() string {
	return string(h.Get())
}

// Get возвращает текущее состояние статистики в формате JSON.
func (h *StatsHolderMock) Get() []byte {
	if stats := h.stats.Load(); stats != nil {
		if raw, ok := stats.([]byte); ok {
			return raw
		}
	}

	return nil
}

// MakeMessageMock создает новое сообщение KAFKA для тестов.
func MakeMessageMock(message *kafka.Message) (*confluent.Message, error) {
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
		Headers: MakeHeadersMock(message.Headers),
	}, nil
}

// ReadMessageMock возвращает представление сообщения KAFKA для тестов.
func ReadMessageMock(message *confluent.Message) *kafka.Message {
	newMessage := &kafka.Message{
		Partition: &message.TopicPartition.Partition,
		Headers:   ReadHeadersMock(message),
	}

	return newMessage.
		WithValue(kafka.ValueBytes(message.Value)).
		WithKeyBytes(message.Key)
}

// EventsPublishMock публикует указанный канал сообщений клиента KAFKA на заданных подписчиках.
func EventsPublishMock(ctx context.Context, events <-chan confluent.Event, subscribers ...bus.Subscriber[confluent.Event]) {
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

// MakeHeadersMock создает новый заголовок сообщения KAFKA для тестов.
func MakeHeadersMock(source collections.MultiMap[string, []byte]) []confluent.Header {
	var headers []confluent.Header
	for key, value := range source {
		for _, v := range value {
			headers = append(headers, confluent.Header{Key: key, Value: v})
		}
	}

	return headers
}

// ReadHeadersMock возвращает представление заголовка сообщения KAFKA для тестов.
func ReadHeadersMock(message *confluent.Message) collections.MultiMap[string, []byte] {
	headers := make(collections.MultiMap[string, []byte])
	for _, header := range message.Headers {
		headers.Append(header.Key, header.Value)
	}

	return headers
}
