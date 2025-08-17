package bus

import (
	"context"
	"math"

	"github.com/wal1251/pkg/core"
	"github.com/wal1251/pkg/tools/collections"
)

var (
	_ Publisher[any]            = (PublisherFn[any])(nil)
	_ Subscriber[any]           = (SubscriberFn[any])(nil)
	_ SubscriberMiddleware[any] = (SubscriberMiddlewareFn[any])(nil)
	_ Subscriber[any]           = SubscriberAdapter[any, int]{}
	_ Subscriber[any]           = SubscriberDemandWrapper[any]{}
)

type (
	// Publisher обобщенный источник событий, для приведения источников к каноничному виду.
	// Термин Publisher использован не корректно. Здесь больше подходит термин EventGenerator.
	Publisher[E any] interface {
		// Subscribe подписывает получатели сообщений (Subscriber) на события источника.
		// При появлении новых событий происходит уведомление получателя.
		// Если подписка на события по каким-то причинам невозможна, метод должен вернуть ошибку.
		Subscribe(ctx context.Context, subscriber Subscriber[E]) error
	}

	// PublisherFn функциональное представление обобщенного источника, семантика функции эквивалентна Publisher.
	PublisherFn[E any] func(ctx context.Context, subscriber Subscriber[E]) error

	// Subscriber обобщенный получатель событий, для приведения получателей к каноничному виду.
	// Термин Subscriber не совсем удачный. Здесь больше подходит EventReceiver.
	Subscriber[E any] interface {
		// Publish передает указанные события в обработку получателю.
		// Если получение событий по каким-то причинам невозможна, метод вернет ошибку.
		Publish(ctx context.Context, events ...E) error

		// Demand возвращает количество событий, которое может быть получено. Если вызов вернет значение
		// 0, то источнику сообщений следует притормозить передачу новых сообщений до тех пор, пока очередной
		// вызов вернет значение отличное от 0. Результат вызова указывает на объем сообщений, который получатель готов
		// принять.
		Demand() int
	}

	// SubscriberFn функциональное представление обобщенного получателя, семантика функции эквивалентна Subscriber.
	SubscriberFn[E any] func(ctx context.Context, events ...E) error

	// SubscriberMiddleware обобщенный встраиваемый посредник получателя.
	SubscriberMiddleware[E any] interface {
		// Call выполняет вызов посредника, в качестве next передается Subscriber, который будет вызван посредником.
		Call(ctx context.Context, events []E, next Subscriber[E]) error
	}

	// SubscriberMiddlewareFn функциональное представление SubscriberMiddleware.
	SubscriberMiddlewareFn[E any] func(ctx context.Context, events []E, next Subscriber[E]) error

	// SubscriberAdapter позволяет преобразовать получатель типа K к получателю типа T, с помощью функции преобразования
	// Transform. Итоговый получатель будет вызывать Subscriber.
	SubscriberAdapter[T, K any] struct {
		Subscriber Subscriber[K]  // Целевой подписчик.
		Transform  core.Map[T, K] // Функция преобразования аргумента подписчика.
	}

	// SubscriberDemandWrapper обертка для получателя типа T, которая позволит задать произвольную функцию OnDemand, для
	// исходного получателя Subscriber. Функция определит поведение метода Demand, который отображает текущую потребность
	// получателя в сообщениях.
	// Может пригодиться, когда есть, к примеру, получатель заданный функцией (для него не определен алгоритм управления
	// индикации потребности), но мы можем дополнительно задать для него этот алгоритм с помощью произвольной функции,
	// обернув в структуру SubscriberDemandWrapper.
	SubscriberDemandWrapper[T any] struct {
		Subscriber Subscriber[T] // Целевой подписчик.
		OnDemand   func() int    // Функция определения потребности получателя в сообщениях.
	}
)

// Publish см. Subscriber.Publish().
func (s SubscriberFn[E]) Publish(ctx context.Context, messages ...E) error {
	if s == nil {
		return nil
	}

	return s(ctx, messages...)
}

// Demand никогда не вернет 0, публикация без ограничений.
func (s SubscriberFn[E]) Demand() int {
	return math.MaxInt
}

// Call см. SubscriberMiddleware.Call().
func (s SubscriberMiddlewareFn[E]) Call(ctx context.Context, messages []E, next Subscriber[E]) error {
	if s == nil {
		return next.Publish(ctx, messages...)
	}

	return s(ctx, messages, next)
}

// Subscribe см. Publisher.Subscribe().
func (p PublisherFn[E]) Subscribe(ctx context.Context, subscriber Subscriber[E]) error {
	if p == nil {
		return nil
	}

	return p(ctx, subscriber)
}

// SubscribeAll возвращает новый SubscriberFn, который последовательно публикует события на всех подписчиках в порядке
// указания в аргументах. Если хоть один Subscriber вернет ошибку, публикация прерывается, будет возвращена ошибка.
func SubscribeAll[E any](subscribers ...Subscriber[E]) SubscriberFn[E] {
	return SubscribeAllWithErr(nil, subscribers...)
}

// SubscribeAllWithErr возвращает новый SubscriberFn, который последовательно публикует события на всех подписчиках в
// порядке указания в аргументах. Возникающие ошибки перехватываются ErrorCallback, он же определяет, будет ли прервана
// цепочка вызовов при ошибке.
func SubscribeAllWithErr[E any](onError core.ErrorCallback, subscribers ...Subscriber[E]) SubscriberFn[E] {
	return func(ctx context.Context, events ...E) error {
		for _, subscriber := range subscribers {
			if err := core.ErrIntercept(subscriber.Publish(ctx, events...), onError); err != nil {
				return err
			}
		}

		return nil
	}
}

// SubscriberWith возвращает новую реализацию Subscriber - подписчик subscriber обернутый в вызовы SubscriberMiddleware.
// Порядок вызова - в порядке передачи в аргументах.
func SubscriberWith[E any](subscriber Subscriber[E], middlewares ...SubscriberMiddleware[E]) Subscriber[E] {
	getNext := func(middleware SubscriberMiddleware[E], next Subscriber[E]) SubscriberFn[E] {
		return func(ctx context.Context, messages ...E) error {
			return middleware.Call(ctx, messages, next)
		}
	}

	next := subscriber
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = getNext(middlewares[i], next)
	}

	return SubscriberDemandWrapper[E]{Subscriber: next, OnDemand: subscriber.Demand}
}

// Publish см. Subscriber.Publish().
func (s SubscriberAdapter[T, K]) Publish(ctx context.Context, events ...T) error {
	return s.Subscriber.Publish(ctx, collections.Map(events, s.Transform.Map)...)
}

// Demand см. Subscriber.Demand().
func (s SubscriberAdapter[T, K]) Demand() int {
	return s.Subscriber.Demand()
}

// Publish см. Subscriber.Publish().
func (s SubscriberDemandWrapper[T]) Publish(ctx context.Context, events ...T) error {
	return s.Subscriber.Publish(ctx, events...)
}

// Demand см. Subscriber.Demand().
func (s SubscriberDemandWrapper[T]) Demand() int {
	if s.OnDemand == nil {
		return math.MaxInt
	}

	return s.OnDemand()
}
