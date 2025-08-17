package bus

import (
	"context"
	"sync/atomic"

	"github.com/wal1251/pkg/tools/collections"
)

const defaultPrefetchFactor = 1.5

var _ Subscriber[any] = (*AsyncDemandingSubscriber[any])(nil)

type (
	// AsyncDemandingSubscriber реализует Subscriber, который потребляет опубликованные сообщения в выделенной go рутине.
	// Опубликованные на подписчике события добавляются в буфер, остаточная емкость этого буфера можно узнать вызовом
	// Demand(), если метод вернет значение 0, это означает что буфер событий заполнен, издателю следует приостановить
	// публикацию событий до тех пор, пока значение не станет больше нуля.
	AsyncDemandingSubscriber[E any] struct {
		size        int
		unprocessed int32
		messages    chan E
		subscriber  Subscriber[E]
		onError     func(err error)
	}
)

// Publish см. Subscriber.Publish().
func (b *AsyncDemandingSubscriber[E]) Publish(_ context.Context, messages ...E) error {
	collections.ForEach(messages, func(msg E) {
		atomic.AddInt32(&b.unprocessed, 1)

		b.messages <- msg
	})

	return nil
}

// Demand см. Subscriber.Demand().
func (b *AsyncDemandingSubscriber[E]) Demand() int {
	demand := b.size - int(atomic.LoadInt32(&b.unprocessed))
	if demand < 0 {
		return 0
	}

	return demand
}

func (b *AsyncDemandingSubscriber[E]) run(ctx context.Context) {
	accept := func(msg E) {
		defer atomic.AddInt32(&b.unprocessed, -1)

		if err := b.subscriber.Publish(ctx, msg); err != nil {
			b.onError(err)
		}
	}

	go func() {
		for msg := range b.messages {
			accept(msg)
		}
	}()
}

// NewAsyncDemandingSubscriber вернет новый экземпляр AsyncDemandingSubscriber, реализующий AsyncDemandingSubscriber.
// При публикации события, оно публикуется для обработки на целевой subscriber. Может быть задан размер буфера событий
// prefetch. Фактически создается буфер в 1.5 больший, на случай, если издатель не сразу прекратит публикацию событий.
// В случае ошибки целевого подписчика будет вызван onError с ошибкой подписчика.
func NewAsyncDemandingSubscriber[E any](ctx context.Context, subscriber Subscriber[E], prefetch int, onError func(err error)) *AsyncDemandingSubscriber[E] {
	if onError == nil {
		onError = func(error) {}
	}

	pub := &AsyncDemandingSubscriber[E]{
		messages:   make(chan E, int(float32(prefetch)*defaultPrefetchFactor)+1),
		subscriber: subscriber,
		size:       prefetch,
		onError:    onError,
	}
	pub.run(ctx)

	return pub
}
