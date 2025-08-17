package bus

import (
	"context"
	"sync"

	"golang.org/x/exp/slices"

	"github.com/wal1251/pkg/tools/collections"
)

var (
	_ EventBus[any] = (*SyncEventBus[any])(nil)
	_ EventBus[any] = (*EventBusAdapter[any, int])(nil)
)

type (
	// EventBus предоставляет абстракцию событийной шины приложения для обмена сообщениям различными компонентами. Шина
	// позволяет публиковать события для различных подписчиков, объединенных в группы - топики. Подписчик должен быть
	// оповещен о событиях публикуемых в топике, на который он подписан.
	EventBus[E any] interface {
		// Notify оповещает подписчиков о новом событии в топике.
		Notify(ctx context.Context, topic string, events ...E) error

		// Subscribe подписывает на события топика указанного подписчика.
		Subscribe(ctx context.Context, topic string, subscriber Subscriber[E]) error

		// Close закрывает шину и освобождает занятые ресурсы.
		Close(ctx context.Context)
	}

	// SyncEventBus реализация событийной шины EventBus по умолчанию, при публикации каждого события синхронно вызывает
	// зарегистрированных подписчиков, если какой-либо из подписчиков вернет ошибку, оповещение прекращается, метод
	// вернет ошибку.
	SyncEventBus[E any] struct {
		lock        sync.RWMutex
		subscribers map[string][]Subscriber[E]
	}

	// EventBusAdapter позволяет преобразовать EventBus типа K в EventBus типа T.
	EventBusAdapter[T, K any] struct {
		// Target целевой EventBus.
		Target EventBus[K]
		// Transform преобразование отправляемых событий.
		Transform func(topic string, event T) (K, error)
		// Read преобразование получаемых событий.
		Read func(ctx context.Context, topic string, event K) (T, error)
		// Обратный вызов в случае ошибки преобразования.
		OnReadError func(ctx context.Context, topic string, event K, err error) error
	}
)

// Notify см. EventBus.Notify().
func (s *SyncEventBus[E]) Notify(ctx context.Context, topic string, events ...E) error {
	for _, subscriber := range s.getSubscribers(topic) {
		if err := subscriber.Publish(ctx, events...); err != nil {
			return err
		}
	}

	return nil
}

// Subscribe см. EventBus.Subscribe().
func (s *SyncEventBus[E]) Subscribe(_ context.Context, topic string, subscriber Subscriber[E]) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.subscribers[topic] = append(s.subscribers[topic], subscriber)

	return nil
}

// Close см. EventBus.Close().
func (s *SyncEventBus[E]) Close(context.Context) {
}

func (s *SyncEventBus[E]) getSubscribers(topic string) []Subscriber[E] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return slices.Clone(s.subscribers[topic])
}

// EventBusSubscriber возвращает нового подписчика, который оповещает о новом событии в указанном топике. Может быть
// использован как адаптер для Subscriber.
func EventBusSubscriber[E any](bus EventBus[E], topic string) SubscriberFn[E] {
	return func(ctx context.Context, events ...E) error {
		return bus.Notify(ctx, topic, events...)
	}
}

// EventBusPublisher возвращает нового издателя, который позволит подписаться на события шины для указанного топика.
// Может быть использован как адаптер для Publisher.
func EventBusPublisher[E any](bus EventBus[E], topic string) PublisherFn[E] {
	return func(ctx context.Context, subscriber Subscriber[E]) error {
		return bus.Subscribe(ctx, topic, subscriber)
	}
}

// NewSyncEventBus вернет новый экземпляр SyncEventBus.
func NewSyncEventBus[E any]() *SyncEventBus[E] {
	return &SyncEventBus[E]{
		subscribers: make(map[string][]Subscriber[E]),
	}
}

// Notify см. EventBus.Notify().
func (a *EventBusAdapter[T, K]) Notify(ctx context.Context, topic string, events ...T) error {
	transformed, err := collections.MapWithErr(events,
		func(event T) (K, error) { return a.Transform(topic, event) })
	if err != nil {
		return err
	}

	return a.Target.Notify(ctx, topic, transformed...)
}

// Subscribe см. EventBus.Subscribe().
func (a *EventBusAdapter[T, K]) Subscribe(ctx context.Context, topic string, subscriber Subscriber[T]) error {
	return a.Target.Subscribe(ctx, topic, &SubscriberDemandWrapper[K]{
		Subscriber: SubscriberFn[K](func(ctx context.Context, events ...K) error {
			transformedList := make([]T, 0, len(events))
			for _, event := range events {
				transformed, err := a.Read(ctx, topic, event)
				if err != nil {
					if a.OnReadError == nil {
						return err
					}

					if err = a.OnReadError(ctx, topic, event, err); err != nil {
						return err
					}

					continue
				}

				transformedList = append(transformedList, transformed)
			}

			if len(transformedList) == 0 {
				return nil
			}

			return subscriber.Publish(ctx, transformedList...)
		}),
		OnDemand: subscriber.Demand,
	})
}

// Close см. EventBus.Close().
func (a *EventBusAdapter[T, K]) Close(ctx context.Context) {
	a.Target.Close(ctx)
}
