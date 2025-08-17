package doubleconfluent

import (
	"context"
	"fmt"
	"sync/atomic"

	confluent "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/wal1251/pkg/core"
	"github.com/wal1251/pkg/core/bus"
	"github.com/wal1251/pkg/providers/kafka"
)

const (
	ConsumerStateIdle     = iota // Состояние потребителя: Покой.
	ConsumerStateActive          // Состояние потребителя: Активен.
	ConsumerStateDisposed        // Состояние потребителя: Остановлен.
)

// ConsumerMock представляет собой мок для клиента-потребителя KAFKA.
type ConsumerMock struct {
	core.Component
	isClosed  chan any
	topics    []string
	state     int32
	consumer  *confluent.Consumer // В реальности это может быть любая реализация, соответствующая интерфейсу
	newPoller func(consumer *confluent.Consumer, subscriber bus.Subscriber[*kafka.Message]) kafka.Poller
}

// NewConsumerMock создает новый мок клиента-потребителя.
func NewConsumerMock(topics []string) *ConsumerMock {
	return &ConsumerMock{
		isClosed: make(chan any),
		topics:   topics,
		state:    ConsumerStateIdle,
	}
}

// Subscribe эмулирует подписку на темы.
func (c *ConsumerMock) Subscribe(_ context.Context, _ bus.Subscriber[*kafka.Message]) error {
	if !atomic.CompareAndSwapInt32(&c.state, ConsumerStateIdle, ConsumerStateActive) {
		state := atomic.LoadInt32(&c.state)

		return fmt.Errorf("%w: state: %d", kafka.ErrInvalidConsumerState, state)
	}

	// Эмулируем успешную подписку
	fmt.Println("Consumer subscribed to topics:", c.topics) //nolint:forbidigo // mock method

	return nil
}

// Close эмулирует закрытие клиента-потребителя.
func (c *ConsumerMock) Close(ctx context.Context) {
	c.dispose(ctx)
	<-c.isClosed
}

// IsActive возвращает true, если клиент активен.
func (c *ConsumerMock) IsActive() bool {
	return atomic.LoadInt32(&c.state) == ConsumerStateActive
}

// Client возвращает клиента более низкого уровня.
func (c *ConsumerMock) Client() *confluent.Consumer {
	return c.consumer
}

// ConfigurePoller устанавливает конструктор для нового поллера.
func (c *ConsumerMock) ConfigurePoller(newPoller func(consumer *confluent.Consumer, subscriber bus.Subscriber[*kafka.Message]) kafka.Poller) {
	c.newPoller = newPoller
}

// dispose освобождает ресурсы.
func (c *ConsumerMock) dispose(_ context.Context) {
	if !atomic.CompareAndSwapInt32(&c.state, ConsumerStateIdle, ConsumerStateDisposed) &&
		!atomic.CompareAndSwapInt32(&c.state, ConsumerStateActive, ConsumerStateDisposed) {
		return
	}

	// Эмулируем логику освобождения ресурсов
	close(c.isClosed)
}

// Receive обрабатывает полученные сообщения.
func (c *ConsumerMock) Receive(_ context.Context, messages ...*kafka.Message) error {
	for _, message := range messages {
		fmt.Println("Received message:", string(message.Value.Must())) //nolint:forbidigo // mock method
	}

	return nil
}
