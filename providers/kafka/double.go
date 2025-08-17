package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	confluent "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/wal1251/pkg/core"
	"github.com/wal1251/pkg/core/bus"
)

var (
	_ Producer = (*ProducerTestDouble)(nil)
	_ Consumer = (*ConsumerTestDouble)(nil)
	_ Poller   = (*PollerMock)(nil)
)

type (
	// ProducerTestDouble эмулирует работу продюсера KAFKA, записывая сообщения в память.
	ProducerTestDouble struct {
		messages []*Message
		mu       sync.Mutex
	}
)

// NewProducerTestDouble создает новый тестовый продюсер.
func NewProducerTestDouble() *ProducerTestDouble {
	return &ProducerTestDouble{
		messages: make([]*Message, 0),
	}
}

// Send эмулирует асинхронную отправку сообщения, записывая его в память.
func (p *ProducerTestDouble) Send(ctx context.Context, message *Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.messages = append(p.messages, message)

	if message.OnAck != nil {
		message.OnAck(ctx)
	}

	return nil
}

// SendSync эмулирует синхронную отправку сообщения, записывая его в память.
func (p *ProducerTestDouble) SendSync(ctx context.Context, message *Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.messages = append(p.messages, message)

	if message.OnAck != nil {
		message.OnAck(ctx)
	}

	return nil
}

// Close эмулирует закрытие продюсера.
func (p *ProducerTestDouble) Close(_ context.Context) {
	// Здесь нет необходимости что-то чистить, так как всё в памяти.
}

// GetMessages возвращает все отправленные сообщения для проверки в тестах.
func (p *ProducerTestDouble) GetMessages() []*Message {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.messages
}

// ClearMessages очищает список сообщений.
func (p *ProducerTestDouble) ClearMessages() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.messages = make([]*Message, 0)
}

// String возвращает строковое представление продюсера для выполнения интерфейса Producer.
func (p *ProducerTestDouble) String() string {
	return fmt.Sprintf("ProducerTestDouble: %d messages", len(p.messages))
}

// ConsumerTestDouble реализация тестового дублера Consumer.
type ConsumerTestDouble struct {
	messages    []*Message                 // Очередь сообщений
	subscribers []bus.Subscriber[*Message] // Список подписчиков
	mu          sync.Mutex                 // Для безопасной работы с подписчиками и сообщениями
	closed      bool                       // Флаг для закрытия потребителя
}

// NewConsumerTestDouble создает новый тестовый дублер Consumer.
func NewConsumerTestDouble() *ConsumerTestDouble {
	return &ConsumerTestDouble{
		messages:    make([]*Message, 0),
		subscribers: make([]bus.Subscriber[*Message], 0),
	}
}

// Subscribe подписывает подписчика на события.
func (c *ConsumerTestDouble) Subscribe(_ context.Context, subscriber bus.Subscriber[*Message]) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return ErrInvalidConsumerState
	}

	c.subscribers = append(c.subscribers, subscriber)

	return nil
}

// Publish публикует сообщения подписчикам.
func (c *ConsumerTestDouble) Publish(ctx context.Context, events ...*Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return ErrInvalidConsumerState
	}

	// Рассылаем сообщения всем подписчикам
	for _, subscriber := range c.subscribers {
		if err := subscriber.Publish(ctx, events...); err != nil {
			return err
		}
	}

	return nil
}

// Demand возвращает максимальное количество сообщений для публикации.
func (c *ConsumerTestDouble) Demand() int {
	// В тестах можем вернуть MaxInt, чтобы не ограничивать публикацию
	return int(^uint(0) >> 1) // Максимальное значение int
}

// Close закрывает потребителя и освобождает ресурсы.
func (c *ConsumerTestDouble) Close(_ context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	c.messages = nil
	c.subscribers = nil
}

// String реализует fmt.Stringer.
func (c *ConsumerTestDouble) String() string {
	return "ConsumerTestDouble"
}

type PollerMock struct {
	lock        sync.RWMutex
	pausedFlag  bool
	messages    []*confluent.Message
	onError     core.ErrorCallback
	pollTimeout time.Duration
}

func NewPollerMock(messages []*confluent.Message, pollTimeout time.Duration, onError core.ErrorCallback) *PollerMock {
	return &PollerMock{
		messages:    messages,
		pollTimeout: pollTimeout,
		onError:     onError,
	}
}

func (p *PollerMock) Poll(_ context.Context) error {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if p.pausedFlag {
		return nil // Если на паузе, просто возвращаем без обработки
	}

	// Симулируем задержку поллинга
	time.Sleep(p.pollTimeout)

	// Обрабатываем первое сообщение, если есть
	if len(p.messages) > 0 {
		msg := p.messages[0]
		p.messages = p.messages[1:]
		// Имитация публикации сообщения подписчику
		fmt.Println("Processed message:", string(msg.Value)) //nolint:forbidigo // mock method
	}

	return nil
}

func (p *PollerMock) Pause(_ context.Context) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.pausedFlag = true

	return nil
}

func (p *PollerMock) Resume(_ context.Context) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.pausedFlag = false

	return nil
}

func (p *PollerMock) IsPaused() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.pausedFlag
}
