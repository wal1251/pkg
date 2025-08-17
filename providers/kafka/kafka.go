// Package kafka предоставляет интерфейсы и их реализации для работы с брокером KAFKA.
package kafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/wal1251/pkg/core/bus"
	"github.com/wal1251/pkg/core/logs"
)

var (
	ErrCancelled            = errors.New("cancelled")
	ErrInvalidConsumerState = errors.New("invalid consumer state")
	ErrMessageNotSent       = errors.New("message not sent")
	ErrTopicCreateFailed    = errors.New("topics creation failed")
	ErrTopicDeleteFailed    = errors.New("topics deletion failed")
)

type (
	// Producer это интерфейс клиента-производителя KAFKA. Предоставляет возможность синхронной и асинхронной отправки
	// сообщений брокеру.
	Producer interface {
		fmt.Stringer

		// Send отправляет сообщение брокеру асинхронно, не ожидает подтверждения о доставке.
		Send(ctx context.Context, messages *Message) error

		// SendSync отправляет сообщение брокеру синхронно. Гарантирует, что брокер получил сообщение.
		SendSync(ctx context.Context, message *Message) error

		// Close закрывает соединение с брокером и высвобождает ресурсы.
		Close(ctx context.Context)
	}

	// Consumer это интерфейс клиента-потребителя KAFKA. Предполагает, что потребитель уже подписан на определенные
	// топики, при подписывании на него публикует сообщения только из этих топиков.
	Consumer interface {
		bus.Publisher[*Message]
		fmt.Stringer

		// Close закрывает соединение с брокером и высвобождает ресурсы.
		Close(ctx context.Context)
	}

	// Poller представляет собой компонент, ответственный за стратегию опроса брокера и подтверждение полученных
	// сообщений (фиксацию смещений).
	Poller interface {
		Poll(ctx context.Context) error
		Pause(ctx context.Context) error
		Resume(ctx context.Context) error
		IsPaused() bool
	}
)

// ProducerToSubscriber приводит Producer к реализации bus.Subscriber. Представление производителя как подписчика, для
// публикации сообщений на брокере.
func ProducerToSubscriber(producer Producer, sync bool) bus.SubscriberFn[*Message] {
	return func(ctx context.Context, messages ...*Message) error {
		return Send(ctx, producer, sync, messages...)
	}
}

// Send отправляет указанные сообщения с помощью клиента Producer, если хоть одно сообщение не удалось отправить вернет
// ошибку ErrMessageNotSent.
func Send(ctx context.Context, producer Producer, sync bool, messages ...*Message) error {
	logger := logs.FromContext(ctx)

	hasErrors := false

	messageSender := producer.Send
	if sync {
		messageSender = producer.SendSync
	}

	for _, message := range messages {
		if err := messageSender(ctx, message); err != nil {
			hasErrors = true

			logger.Err(err).Msgf("failed to send message: %v", producer)
		}
	}

	if hasErrors {
		return ErrMessageNotSent
	}

	return nil
}
