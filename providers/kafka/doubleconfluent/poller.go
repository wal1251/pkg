package doubleconfluent

import (
	"context"
	"sync"

	confluent "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// MockPoller - моковая реализация kafka.Poller.
type MockPoller struct {
	mu           sync.Mutex
	messages     []*confluent.Message // Сохраненные сообщения для проверки
	paused       bool                 // Флаг приостановки
	errOnPoll    error                // Ошибка, возвращаемая при опросе
	errOnPublish error                // Ошибка, возвращаемая при публикации
}

// NewMockPoller создает новый MockPoller с заданными ошибками.
func NewMockPoller(errOnPoll, errOnPublish error) *MockPoller {
	return &MockPoller{
		messages:     make([]*confluent.Message, 0),
		errOnPoll:    errOnPoll,
		errOnPublish: errOnPublish,
	}
}

// Poll имитирует опрос брокера на предмет наличия новых сообщений.
func (mp *MockPoller) Poll(_ context.Context) error {
	if mp.errOnPoll != nil {
		return mp.errOnPoll // Возвращаем ошибку, если она задана
	}

	// Имитация получения сообщения (можно добавить логику для генерации тестовых сообщений)
	msg := &confluent.Message{
		// Заполните поля сообщения для тестов
	}

	mp.messages = append(mp.messages, msg)

	if mp.paused {
		return nil // Если приостановлено, не публикуем сообщения
	}

	// Имитация публикации сообщения (проверка ошибки)
	if mp.errOnPublish != nil {
		return mp.errOnPublish // Возвращаем ошибку, если она задана
	}

	// Здесь можно добавить логику для вызова подписчика

	return nil
}

// IsPaused возвращает true, если публикация сообщений приостановлена.
func (mp *MockPoller) IsPaused() bool {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	return mp.paused
}

// Pause приостанавливает публикацию сообщений.
func (mp *MockPoller) Pause(_ context.Context) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.paused = true

	return nil
}

// Resume возобновляет публикацию сообщений.
func (mp *MockPoller) Resume(_ context.Context) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.paused = false

	return nil
}

// GetMessages возвращает сохраненные сообщения.
func (mp *MockPoller) GetMessages() []*confluent.Message {
	return mp.messages
}

// ClearMessages очищает сохраненные сообщения.
func (mp *MockPoller) ClearMessages() {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.messages = nil
}
