package doubleconfluent

import (
	"context"

	"github.com/wal1251/pkg/providers/kafka/confluent"

	"github.com/wal1251/pkg/providers/kafka"
)

// MockProducer - моковая реализация клиента-производителя KAFKA.
type MockProducer struct {
	messages  []kafka.Message // Сохраненные сообщения для проверки
	errToSend error           // Ошибка, возвращаемая при отправке
}

// NewMockProducer создает новый MockProducer с заданной ошибкой при отправке.
func NewMockProducer(errToSend error) *MockProducer {
	return &MockProducer{
		messages:  make([]kafka.Message, 0),
		errToSend: errToSend,
	}
}

// Send асинхронная отправка сообщений брокеру.
func (mp *MockProducer) Send(ctx context.Context, message *kafka.Message) error {
	return mp.send(ctx, message, false)
}

// SendSync синхронная отправка сообщений брокеру.
func (mp *MockProducer) SendSync(ctx context.Context, message *kafka.Message) error {
	return mp.send(ctx, message, true)
}

// Close освобождает ресурс.
func (mp *MockProducer) Close(_ context.Context) {
	// todo: добавить логику для освобождения ресурсов, если необходимо
}

// Client возвращает клиент более низкого уровня (nil для моков).
func (mp *MockProducer) Client() *confluent.Producer {
	return nil
}

// send - основная логика отправки сообщения (мок).
func (mp *MockProducer) send(_ context.Context, message *kafka.Message, _ bool) error {
	if mp.errToSend != nil {
		return mp.errToSend // Возвращаем ошибку, если она задана
	}

	// Сохраняем сообщение для проверки
	mp.messages = append(mp.messages, *message)

	return nil
}

// GetMessages возвращает сохраненные сообщения.
func (mp *MockProducer) GetMessages() []kafka.Message {
	return mp.messages
}
