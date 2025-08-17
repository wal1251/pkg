package doubleconfluent // Не собираем этот пример, т.к. пакет требует внешней зависимости.

import (
	"context"
	"fmt"
	"log"
	"sync"

	confluent "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/wal1251/pkg/core/bus"
	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/providers/kafka"
)

func ExampleProducer_Send() {
	// Создадим контекст с логером.
	logger := logs.Logger(&logs.Config{Level: "trace", Pretty: true})
	ctx := logger.WithContext(context.TODO())

	// Создадим клиента-производителя сообщений.
	producer := NewMockProducer(nil)
	events := make(chan confluent.Event)

	// Будем логировать события производителя.
	EventsPublishMock(ctx, events)

	// Создаем мок-консюмер
	consumer := NewConsumerMock([]string{"my.topic.test01"})
	defer consumer.Close(ctx)
	if err := consumer.Subscribe(ctx, bus.SubscriberFn[*kafka.Message](func(ctx context.Context, messages ...*kafka.Message) error {
		for _, message := range messages {
			fmt.Println("Consumed message:", string(message.Value.Must()))
		}
		return nil
	})); err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	// Запускаем горутину для получения событий из канала
	wg.Add(1)
	go func() {
		defer wg.Done()
		for event := range events {
			if messageEvent, ok := event.(*MessageEvent); ok {
				if err := consumer.Receive(ctx, messageEvent.Message); err != nil {
					log.Println("Error receiving message:", err)
				}
			}
		}
	}()

	for i := 1; i <= 2; i++ {
		// Отправляем сообщение синхронно в топик my.topic.test01.
		message := (&kafka.Message{Topic: "my.topic.test01"}).
			WithValue(kafka.ValueBytes([]byte(fmt.Sprintf("hello %d!", i))))
		fmt.Println("Sending message", i, ":", string(message.Value.Must()))
		if err := producer.SendSync(ctx, message); err != nil {
			log.Fatal(err)
		}
		events <- &MessageEvent{Message: message} // Публикуем событие в канал
	}

	close(events)
	wg.Wait()
}

func ExampleEventBusDouble() {
	ctx := context.TODO()

	// Создаем тестовый EventBus.
	eventBus := NewEventBusDouble("test.prefix.", true)
	defer eventBus.Close(ctx)

	// Создаем подписчика на топик "my.topic.test01".
	err := eventBus.Subscribe(ctx, "my.topic.test01", bus.SubscriberFn[*kafka.Message](func(ctx context.Context, messages ...*kafka.Message) error {
		for _, message := range messages {
			fmt.Println("Consumed message:", string(message.Value.Must()))
		}
		return nil
	}))
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// Запускаем горутину для имитации асинхронной обработки сообщений.
	go func() {
		defer wg.Done()
		for i := 1; i <= 2; i++ {
			// Создаем сообщение для топика "my.topic.test01".
			message := &kafka.Message{
				Topic: "my.topic.test01",
				Value: kafka.ValueBytes([]byte(fmt.Sprintf("hello %d!", i))),
			}
			fmt.Println("Sending message", i, ":", string(message.Value.Must()))

			// Отправляем сообщение через EventBusTestDouble.
			if err = eventBus.Notify(ctx, "", message); err != nil {
				log.Fatal(err)
			}
		}
	}()

	wg.Wait()
}
