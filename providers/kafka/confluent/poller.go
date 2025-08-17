package confluent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	confluent "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/wal1251/pkg/core"
	"github.com/wal1251/pkg/core/bus"
	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/providers/kafka"
)

const (
	PollerComponentName = "confluent-kafka-go-poller" // Имя компонента
)

var (
	_ kafka.Poller = (*Poller)(nil)
	_ kafka.Poller = (*PollerRespectingDemand)(nil)
)

var ErrSubscriberPanics = errors.New("subscriber panics")

type (
	// Poller предоставляет реализацию kafka.Poller. Реализует стратегии опроса брокера KAFKA.
	Poller struct {
		core.Component

		lock       sync.RWMutex
		pausedFlag bool

		consumer    *confluent.Consumer
		onError     core.ErrorCallback
		pollTimeout time.Duration
		subscriber  bus.Subscriber[*confluent.Message]
	}

	// PollerRespectingDemand предоставляет реализацию kafka.Poller, который учитывает требование подписчика о количестве
	// производимых сообщений. Если потребность подписчика будет равна 0, то при очередном вызове метода
	// PollerRespectingDemand.Poll() потребить KAfKA будет приостановлен - новых сообщений публиковаться не будет.
	PollerRespectingDemand struct {
		*Poller
		suspend func() bool
	}
)

// Poll опрос брокера на предмет наличия новых сообщений.
func (p *Poller) Poll(ctx context.Context) error {
	logger := logs.FromContext(ctx)

	event := p.consumer.Poll(kafka.Milliseconds(p.pollTimeout))
	if event != nil {
		if err := EventsLogger(p)(ctx, event); err != nil {
			logger.Err(err).Msg("error while logging consumer events")
		}
	}

	switch typedEvent := event.(type) {
	case *confluent.Message:
		if typedEvent.TopicPartition.Error != nil {
			return core.ErrIntercept(typedEvent.TopicPartition.Error, p.onError)
		}

		if p.IsPaused() {
			logger.Warn().Msg("message income on paused consumer")

			// TODO: Нужно разобраться, почему можем сюда попасть и найти более элегантное решение.
			if err := p.Pause(ctx); err != nil {
				return err
			}
		}

		return p.subscriber.Publish(ctx, typedEvent)

	case confluent.AssignedPartitions:
		logger.Info().Msgf("assigned: %v", typedEvent)

		if err := p.consumer.Assign(typedEvent.Partitions); err != nil {
			return core.ErrIntercept(fmt.Errorf("error while consumer %v assigning %v: %w", p.consumer, typedEvent, err), p.onError)
		}

	case confluent.RevokedPartitions:
		logger.Info().Msgf("revoked: %v", typedEvent)

		if err := p.consumer.Unassign(); err != nil {
			return core.ErrIntercept(fmt.Errorf("error while consumer %v unassigning %v: %w", p.consumer, typedEvent, err), p.onError)
		}

	case confluent.Error:
		if typedEvent.IsFatal() {
			return core.ErrIntercept(fmt.Errorf("kafka consumer fatal error %v: %w", p.consumer, typedEvent), p.onError)
		}

		logger.Err(typedEvent).Msg("error occurred while consuming")
	}

	return nil
}

// IsPaused вернет true, если публикация сообщений приостановлена.
func (p *Poller) IsPaused() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.pausedFlag
}

// Pause приостанавливает публикацию сообщений при дальнейших опросах брокера.
func (p *Poller) Pause(ctx context.Context) error {
	logger := logs.FromContext(ctx)

	return p.togglePause(true, func() error {
		assignment, err := p.consumer.Assignment()
		if err != nil {
			return core.ErrIntercept(fmt.Errorf("can't retrieve consumer assignment: %w", err), p.onError)
		}

		if err = p.consumer.Pause(assignment); err != nil {
			return core.ErrIntercept(fmt.Errorf("can't pause consumer: %w", err), p.onError)
		}

		logger.Debug().Msg("consumer paused")

		return nil
	})
}

// Resume возобновляет публикацию сообщений при дальнейших опросах брокера.
func (p *Poller) Resume(ctx context.Context) error {
	logger := logs.FromContext(ctx)

	return p.togglePause(false, func() error {
		assignment, err := p.consumer.Assignment()
		if err != nil {
			return core.ErrIntercept(fmt.Errorf("can't retrieve consumer assignment: %w", err), p.onError)
		}

		if err = p.consumer.Resume(assignment); err != nil {
			return core.ErrIntercept(fmt.Errorf("can't resume consumer: %w", err), p.onError)
		}

		logger.Debug().Msg("consumer resumed")

		return nil
	})
}

func (p *Poller) togglePause(desired bool, toggle func() error) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.pausedFlag == desired {
		return nil
	}

	if err := toggle(); err != nil {
		return err
	}

	p.pausedFlag = desired

	return nil
}

// Poll опрос брокера на предмет наличия новых сообщений. Новые сообщения будут публиковаться в зависимости от
// потребности подписчика.
func (p *PollerRespectingDemand) Poll(ctx context.Context) error {
	if p.suspend() {
		if err := p.Pause(ctx); err != nil {
			return err
		}
	} else {
		if err := p.Resume(ctx); err != nil {
			return err
		}
	}

	if err := p.Poller.Poll(ctx); err != nil {
		return err
	}

	return nil
}

// NewPoller возвращает новый Poller для выполнения опроса брокера. Необходимо передать: клиент-потребитель брокера,
// таймаут опроса и подписчика на сообщения брокера.
func NewPoller(consumer *confluent.Consumer, pollTimeout time.Duration, subscriber bus.Subscriber[*confluent.Message], onError core.ErrorCallback) *Poller {
	return &Poller{
		Component:   core.NewDefaultComponent(PollerComponentName, consumer.String()),
		consumer:    consumer,
		pollTimeout: pollTimeout,
		subscriber:  subscriber,
		onError:     onError,
	}
}

// NewPollerRespectingDemand возвращает новый PollerRespectingDemand для выполнения опроса брокера. Необходимо передать:
// клиент-потребитель брокера, таймаут опроса и подписчика на сообщения брокера.
func NewPollerRespectingDemand(
	consumer *confluent.Consumer,
	pollTimeout time.Duration,
	publisher bus.Subscriber[*confluent.Message],
	onError core.ErrorCallback,
) *PollerRespectingDemand {
	return &PollerRespectingDemand{
		Poller:  NewPoller(consumer, pollTimeout, publisher, onError),
		suspend: func() bool { return publisher.Demand() <= 0 },
	}
}

// PollerDefaultFactory возвращает новую фабрику для конструирования реализаций kafka.Poller по-умолчанию.
func PollerDefaultFactory(cfg *kafka.Config, onError core.ErrorCallback) func(*confluent.Consumer, bus.Subscriber[*kafka.Message]) kafka.Poller {
	return func(consumer *confluent.Consumer, subscriber bus.Subscriber[*kafka.Message]) kafka.Poller {
		return NewPollerRespectingDemand(consumer, cfg.PollTimeout,
			bus.SubscriberWith[*confluent.Message](
				&bus.SubscriberAdapter[*confluent.Message, *kafka.Message]{
					Subscriber: subscriber,
					Transform:  DefaultMessageReader(consumer, onError),
				},
				WithRecover(consumer, true, onError),
			),
			onError,
		)
	}
}

// WithCommitAfterEveryPublish вернет посредника для подписчика на сообщения KAFKA, который выполняет фиксацию смещения
// после обработки сообщения. Для использования в Poller и его производных.
func WithCommitAfterEveryPublish(consumer *confluent.Consumer, onError core.ErrorCallback) bus.SubscriberMiddlewareFn[*confluent.Message] {
	return func(ctx context.Context, messages []*confluent.Message, next bus.Subscriber[*confluent.Message]) error {
		defer func() {
			if err := CommitMessage(consumer, messages...); err != nil {
				logs.FromContext(ctx).Err(err).Msg("commit failed")
				core.ErrNotify(err, onError)
			}
		}()

		return next.Publish(ctx, messages...)
	}
}

// WithRecover вернет посредника для подписчика на сообщения KAFKA, в случае паники выполняет восстановление, логирует
// ошибку и продолжает выполнение подписчика. Для использования в Poller и его производных.
func WithRecover(consumer *confluent.Consumer, mustAck bool, onError core.ErrorCallback) bus.SubscriberMiddlewareFn[*confluent.Message] {
	return func(ctx context.Context, messages []*confluent.Message, next bus.Subscriber[*confluent.Message]) error {
		logger := logs.FromContext(ctx)
		defer func() {
			if x := recover(); x != nil {
				logger.Error().Stack().Msgf("subscriber panic: %v", x)
				core.ErrNotify(fmt.Errorf("%w: %v", ErrSubscriberPanics, x), onError)

				if mustAck {
					if err := CommitMessage(consumer, messages...); err != nil {
						logger.Err(err).Msg("failed to commit message")
					}
				}
			}
		}()

		return next.Publish(ctx, messages...)
	}
}
