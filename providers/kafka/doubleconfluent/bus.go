package doubleconfluent

import (
	"context"
	"sync"

	"github.com/wal1251/pkg/providers/kafka"

	"github.com/wal1251/pkg/core/bus"
)

// EventBusDouble — тестовый двойник для EventBus.
type EventBusDouble struct {
	prefix      string
	messages    map[string][]*kafka.Message                 // Карта для хранения сообщений по топикам
	subscribers map[string][]bus.Subscriber[*kafka.Message] // Карта для подписчиков по топикам
	mu          sync.Mutex                                  // Для безопасной работы с картами
	syncNotify  bool                                        // Для симуляции синхронного уведомления
}

// NewEventBusDouble создает новый тестовый EventBus.
func NewEventBusDouble(prefix string, syncNotify bool) *EventBusDouble {
	return &EventBusDouble{
		prefix:      prefix,
		messages:    make(map[string][]*kafka.Message),
		subscribers: make(map[string][]bus.Subscriber[*kafka.Message]),
		syncNotify:  syncNotify,
	}
}

// Notify выполняет публикацию событий в топике name, если топик - пустая строка, тогда топик для публикации будет взят
// из сообщения. Перед публикацией к имени топика добавляется префикс.
func (b *EventBusDouble) Notify(ctx context.Context, name string, events ...*kafka.Message) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, event := range events {
		topic := name
		if topic == "" {
			topic = event.Topic
		}
		topic = b.prefix + topic

		b.messages[topic] = append(b.messages[topic], event)

		if subscribers, ok := b.subscribers[topic]; ok {
			for _, subscriber := range subscribers {
				if err := subscriber.Publish(ctx, event); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Subscribe регистрирует подписчика топика, к имени топика добавляется префикс. В опубликованных сообщениях топик указан
// без префикса.
func (b *EventBusDouble) Subscribe(_ context.Context, name string, subscriber bus.Subscriber[*kafka.Message]) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	topic := b.prefix + name

	b.subscribers[topic] = append(b.subscribers[topic], subscriber)

	return nil
}

// Close используется для очистки ресурсов.
func (b *EventBusDouble) Close(_ context.Context) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Очищаем карты
	b.messages = make(map[string][]*kafka.Message)
	b.subscribers = make(map[string][]bus.Subscriber[*kafka.Message])
}

// GetMessages возвращает сообщения для заданного топика.
func (b *EventBusDouble) GetMessages(topic string) []*kafka.Message {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.messages[topic]
}

// ClearMessages очищает все сообщения для заданного топика.
func (b *EventBusDouble) ClearMessages(topic string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.messages[topic] = []*kafka.Message{}
}
