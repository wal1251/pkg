package confluent

import (
	"context"
	"fmt"

	confluent "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/wal1251/pkg/core"
	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/providers/kafka"
)

const ProducerComponentName = "confluent-kafka-go-producer" // Имя компонента клиента-производителя.

var _ kafka.Producer = (*Producer)(nil)

type (
	// Producer реализация клиента-потребителя KAFKA, реализует kafka.Producer.
	Producer struct {
		core.Component

		producer *confluent.Producer
	}
)

// Send асинхронная отправка сообщений брокеру.
func (p *Producer) Send(ctx context.Context, message *kafka.Message) error {
	return p.send(ctx, message, false)
}

// SendSync синхронная отправка сообщений брокеру.
func (p *Producer) SendSync(ctx context.Context, message *kafka.Message) error {
	return p.send(ctx, message, true)
}

// Close освобождает ресурс.
func (p *Producer) Close(ctx context.Context) {
	logs.LocalContext(p).To(logs.FromContext(ctx).Debug).Msg("closing producer")

	p.producer.Close()
}

func (p *Producer) String() string {
	return p.producer.String()
}

// Client возвращает клиент более низкого уровня.
func (p *Producer) Client() *confluent.Producer {
	return p.producer
}

func (p *Producer) send(ctx context.Context, message *kafka.Message, isSync bool) error {
	logs.LocalContext(p).To(logs.FromContext(ctx).Debug).Msg("sending message")

	kafkaMessage, err := MakeMessage(message)
	if err != nil {
		return err
	}

	var deliveryChan chan confluent.Event
	if isSync {
		deliveryChan = make(chan confluent.Event)
	}

	if err = p.producer.Produce(kafkaMessage, deliveryChan); err != nil {
		return fmt.Errorf("can't produce kafka message: %w", err)
	}

	if isSync {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: stopped delivery waiting", kafka.ErrCancelled)
		case event := <-deliveryChan:
			if e, ok := event.(confluent.Error); ok {
				return fmt.Errorf("failed message delivery: %w", e)
			}
		}
	}

	return nil
}

// NewProducer возвращает новый Producer, клиент-производитель KAFKA.
func NewProducer(cfg *kafka.Config, opts ...kafka.ClientOption) (*Producer, error) {
	clientConfig, err := kafka.ProducerConfig(cfg, opts...)
	if err != nil {
		return nil, fmt.Errorf("can't create kafka producer: %w", err)
	}

	kafkaProducer, err := confluent.NewProducer(ConfigMap(clientConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	producer := &Producer{
		Component: core.NewDefaultComponent(ProducerComponentName, kafkaProducer.String()),

		producer: kafkaProducer,
	}

	return producer, nil
}
