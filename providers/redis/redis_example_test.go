package redis_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wal1251/pkg/providers/redis"
)

func ExampleClient() {
	ctx := context.Background()

	// Объявляем конфигурацию для клиента Redis.
	//
	// Для загрузки конфигурации из переменных окружения в вашем приложении,
	// можно использовать метод CfgFromViper.
	cfg := &redis.Config{
		Host:               "localhost",
		Port:               "6379",
		MaxBulkRequestSize: 100,
	}

	// Запускаем сервер Redis на основе miniredis.
	tr := redis.NewTestRedisServer()
	if err := tr.Run(*cfg); err != nil {
		log.Fatalf("Failed to start Redis server: %v", err)
	}
	defer tr.Close()

	// Создание экземпляра клиента Redis.
	client, err := redis.NewClient(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create Redis client: %v", err)
	}
	// Не забываем закрыть соединение с Redis после завершения работы.
	defer client.Close(ctx)

	// Пример установки значения с заданным временем жизни.
	err = client.Set(ctx, "key1", "value1", 10*time.Second)
	if err != nil {
		log.Fatalf("Failed to set value in Redis: %v", err)
	}

	// Пример получения значения (строки).
	value, err := client.Get(ctx, "key1")
	if err != nil {
		log.Fatalf("Failed to get value from Redis: %v", err)
	}
	strVal, err := value.String()
	if err != nil {
		log.Fatalf("Failed to get value from Redis: %v", err)
	}
	fmt.Printf("Retrieved string value: %s\n", strVal)

	// Пример получения значения (структуры).
	type Human struct {
		Name string
		Age  int
	}
	john := Human{
		Name: "John",
		Age:  23,
	}

	err = client.Set(ctx, "key2", john, 10*time.Second)
	if err != nil {
		log.Fatalf("Failed to set value in Redis: %v", err)
	}

	value, err = client.Get(ctx, "key2")
	if err != nil {
		log.Fatalf("Failed to get value from Redis: %v", err)
	}

	var johnFromRedis Human
	err = value.Struct(&johnFromRedis)
	if err != nil {
		log.Fatalf("Failed to get value from Redis: %v", err)
	}

	fmt.Printf("Retrieved struct value: %+v\n", johnFromRedis)

	// Пример получения списка значений (где один из ключей не существует).
	keys := []string{"key1", "nonexistentKey"}
	values, err := client.GetList(ctx, keys...)
	if err != nil {
		log.Fatalf("Failed to get list of values from Redis: %v", err)
	}
	for i, v := range values {
		if v != nil {
			strVal, err = v.String()
			if err != nil {
				log.Fatalf("Failed to get value from Redis: %v", err)
			}
			fmt.Printf("Value at index %d: %s\n", i, strVal)
		} else {
			fmt.Printf("Value at index %d is nil\n", i)
		}
	}

	// Пример удаления ключа.
	delCount, err := client.Delete(ctx, "key1")
	if err != nil {
		log.Fatalf("Failed to delete key from Redis: %v", err)
	}
	fmt.Printf("Deleted %d keys\n", delCount)

	// Output:
	// Retrieved string value: value1
	// Retrieved struct value: {Name:John Age:23}
	// Value at index 0: value1
	// Value at index 1 is nil
	// Deleted 1 keys
}
