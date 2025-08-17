//go:build exclude

package confluent_test // Не собираем этот пример, т.к. пакет требует внешней зависимости.

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	vendor "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/wal1251/pkg/core"
	"github.com/wal1251/pkg/core/bus"
	"github.com/wal1251/pkg/core/cfg/viperx"
	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/providers/kafka"
	"github.com/wal1251/pkg/providers/kafka/confluent"
	"github.com/wal1251/pkg/tools/collections"
)

func ExampleProducer_Send() {
	// Создадим контекст с логером.
	logger := logs.Logger(&logs.Config{Level: "trace", Pretty: true})
	ctx := logger.WithContext(context.TODO())

	// Имитация окружения.
	if err := os.Setenv(kafka.CfgKeyHosts.String(), "localhost:29092"); err != nil {
		log.Fatal(err)
	}

	// Прочитаем конфигурацию из переменных окружения.
	cfg := kafka.CfgFromViper(viperx.EnvLoader(""))

	// Создадим клиента-производителя сообщений.
	producer, err := confluent.NewProducer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Будем логировать события производителя.
	confluent.EventsPublish(ctx, producer.Client().Events(), confluent.EventsLogger(nil))

	for i := 1; i <= 10; i++ {
		// Отправляем сообщение синхронно в топик my.topic.test02.
		// Лучше использовать асинхронный вариант: producer.Send().
		if err = producer.SendSync(ctx,
			(&kafka.Message{Topic: "my.topic.test03"}).
				WithValue(kafka.ValueBytes([]byte(fmt.Sprintf("hello %d!", i)))),
		); err != nil {
			log.Fatal(err)
		}
	}
}

func ExampleConsumer_Subscribe() {
	// Создадим контекст с логером.
	logger := logs.Logger(&logs.Config{Level: "trace", Pretty: true})
	ctx, cancel := context.WithCancel(logger.WithContext(context.TODO()))

	// Имитация окружения.
	if err := os.Setenv(kafka.CfgKeyHosts.String(), "localhost:29092"); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv(kafka.CfgKeyGroupID.String(), "test.group.01"); err != nil {
		log.Fatal(err)
	}

	// Прочитаем конфигурацию из переменных окружения.
	cfg := kafka.CfgFromViper(viperx.EnvLoader(""))

	// Создадим клиента-читателя топика my.topic.test02.
	consumer, err := confluent.NewConsumer(cfg,
		collections.Single("my.topic.test02"),
		core.ErrorCallbackFn(func(err error) bool {
			// При ошибке потребителя закроем приложение.
			log.Fatal(err)
			return false
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer consumer.Close(ctx) // Будет ждать закрытия потребителя.
	defer cancel()            // При отмене контекста потребление прекратится, если вызываем закрытие, можно не отменять контекст.

	// Подпишем клиента.
	if err = consumer.Subscribe(ctx, bus.SubscriberFn[*kafka.Message](func(ctx context.Context, messages ...*kafka.Message) error {
		for _, message := range messages {
			// Обязательно подтвердим прочитанное сообщение.
			message.Ack(ctx)

			// Прочитаем тело сообщения.
			fmt.Println(consumer, string(message.Value.Must()))
		}
		return nil
	})); err != nil {
		log.Fatal(err)
	}

	// Подождем некоторое время, потом остановим потребление.
	time.Sleep(30 * time.Second)
}

func ExampleConsumer_ConfigurePoller() {
	// Создадим контекст с логером.
	logger := logs.Logger(&logs.Config{Level: "trace", Pretty: true})
	ctx, cancel := context.WithCancel(logger.WithContext(context.TODO()))

	// Имитация окружения.
	if err := os.Setenv(kafka.CfgKeyHosts.String(), "localhost:29092"); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv(kafka.CfgKeyGroupID.String(), "test.group.01"); err != nil {
		log.Fatal(err)
	}

	// Прочитаем конфигурацию из переменных окружения.
	cfg := kafka.CfgFromViper(viperx.EnvLoader(""))

	errCallback := core.ErrorCallbackFn(func(err error) bool {
		// При ошибке потребителя закроем приложение.
		log.Fatal(err)
		return false
	})

	// Создадим клиента-читателя топика my.topic.test03.
	consumer, err := confluent.NewConsumer(cfg,
		collections.Single("my.topic.test03"),
		errCallback,
	)

	// Сконфигурируем кастомный алгоритм опроса брокера.
	// Теперь не надо подтверждать каждое сообщение вручную, оно будет подтверждено сразу после получения.
	consumer.ConfigurePoller(func(consumer *vendor.Consumer, subscriber bus.Subscriber[*kafka.Message]) kafka.Poller {
		return confluent.NewPoller(consumer, 70*time.Millisecond,
			bus.SubscriberWith[*vendor.Message](
				&bus.SubscriberAdapter[*vendor.Message, *kafka.Message]{
					Subscriber: subscriber,
					Transform:  confluent.DefaultMessageReader(consumer, errCallback),
				},
				confluent.WithCommitAfterEveryPublish(consumer, errCallback),
				confluent.WithRecover(consumer, true, errCallback),
			),
			errCallback,
		)
	})

	if err != nil {
		log.Fatal(err)
	}

	defer consumer.Close(ctx) // Будет ждать закрытия потребителя.
	defer cancel()            // При отмене контекста потребление прекратится, если вызываем закрытие, можно не отменять контекст.

	// Подпишем клиента.
	if err = consumer.Subscribe(ctx, bus.SubscriberFn[*kafka.Message](func(ctx context.Context, messages ...*kafka.Message) error {
		for _, message := range messages {
			// Прочитаем тело сообщения.
			fmt.Println(consumer, string(message.Value.Must()))
		}
		return nil
	})); err != nil {
		log.Fatal(err)
	}

	// Подождем некоторое время, потом остановим потребление.
	time.Sleep(30 * time.Second)
}

func ExampleEventBus() {
	// Создадим контекст с логером.
	logger := logs.Logger(&logs.Config{Level: "trace", Pretty: true})
	ctx, cancel := context.WithCancel(logger.WithContext(context.TODO()))

	// Имитация окружения.
	if err := os.Setenv(kafka.CfgKeyPrefix.String(), "my"); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv(kafka.CfgKeyHosts.String(), "localhost:29092"); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv(kafka.CfgKeyGroupID.String(), "test.group.01"); err != nil {
		log.Fatal(err)
	}

	// Прочитаем конфигурацию из переменных окружения.
	cfg := kafka.CfgFromViper(viperx.EnvLoader(""))

	// Создадим событийную шину.
	eventBus, err := confluent.NewEventBus(ctx,
		cfg,
		core.ErrorCallbackFn(func(err error) bool {
			// При ошибке потребителя закроем приложение.
			log.Fatal(err)
			return false
		}))
	if err != nil {
		log.Fatal(err)
	}

	defer eventBus.Close(ctx) // Будет ждать закрытия потребителя.
	defer cancel()            // При отмене контекста потребление прекратится, если вызываем закрытие, можно не отменять контекст.

	var wg sync.WaitGroup
	wg.Add(10)

	processMessage := func(message *kafka.Message) {
		defer wg.Done()        // Можно больше не дать этого сообщения.
		defer message.Ack(ctx) // Обязательно подтвердим прочитанное сообщение.

		// Прочитаем тело сообщения.
		fmt.Println(string(message.Value.Must()))
	}

	for i := 1; i <= 10; i++ {
		// Отправляем сообщение асинхронно, в топик: my.dev.topic.test03.
		// Префикс my.dev теперь указывать не надо, он задан в конфиге.
		if err = eventBus.Notify(ctx,
			"topic.test03",
			(&kafka.Message{}).WithValue(kafka.ValueBytes([]byte(fmt.Sprintf("hello %d!", i)))),
		); err != nil {
			log.Fatal(err)
		}
	}

	// Подпишемся на тему my.dev.topic.test03 шины. Префикс my.dev указывать не надо, он задан в конфиге.
	if err = eventBus.Subscribe(ctx, "topic.test03", bus.SubscriberFn[*kafka.Message](func(ctx context.Context, messages ...*kafka.Message) error {
		for _, message := range messages {
			processMessage(message)
		}
		return nil
	})); err != nil {
		log.Fatal(err)
	}

	// Подождем считывания десяти сообщений и закроемся.
	wg.Wait()
}

// ExampleEventBus_sync пример создания Kafka шины с синхронным производителем (для гарантии доставки сообщений).
func ExampleEventBus_sync() {
	// Создадим контекст с логером.
	logger := logs.Logger(&logs.Config{Level: "trace", Pretty: true})
	ctx, cancel := context.WithCancel(logger.WithContext(context.TODO()))

	// Устанавливаем конфигурацию для производителя.
	producerCfg := map[string]any{
		kafka.ClientConfigRequestRequiredACKs: kafka.ClientConfigRequestRequiredACKsOne,
	}
	raw, err := json.Marshal(producerCfg)
	if err := os.Setenv(kafka.CfgKeyProducer.String(), string(raw)); err != nil {
		log.Fatal(err)
	}

	// Прочитаем конфигурацию из переменных окружения.
	cfg := kafka.CfgFromViper(viperx.EnvLoader(""))

	// Создадим событийную шину.
	eventBus, err := confluent.NewEventBus(ctx,
		cfg,
		core.ErrorCallbackFn(func(err error) bool {
			// При ошибке потребителя закроем приложение.
			log.Fatal(err)
			return false
		}))
	if err != nil {
		log.Fatal(err)
	}

	defer eventBus.Close(ctx) // Будет ждать закрытия потребителя.
	defer cancel()            // При отмене контекста потребление прекратится, если вызываем закрытие, можно не отменять контекст.
}
