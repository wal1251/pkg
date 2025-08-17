package kafka

import (
	"context"
	"fmt"
	"testing"

	"github.com/wal1251/pkg/core/bus"
)

func TestProducerConsumerSync(t *testing.T) {
	ctx := context.Background()

	producer := NewProducerTestDouble()
	consumer := NewConsumerTestDouble()

	err := consumer.Subscribe(ctx, bus.SubscriberFn[*Message](func(ctx context.Context, messages ...*Message) error {
		for _, message := range messages {
			// Обрабатываем каждое полученное сообщение
			fmt.Println("Received message:", string(message.Value.Must()))
		}
		return nil
	}))
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	message := &Message{
		Topic: "my.topic.test01",
		Value: ValueBytes([]byte("Hello, Kafka!")),
	}

	err = producer.SendSync(ctx, message)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	err = consumer.Publish(ctx, message)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	if len(producer.GetMessages()) == 0 {
		t.Fatal("Expected message to be sent, but got none")
	}

	receivedMessages := producer.GetMessages()
	if len(receivedMessages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(receivedMessages))
	}

	expectedMessage := "Hello, Kafka!"
	if string(receivedMessages[0].Value.Must()) != expectedMessage {
		t.Fatalf("Expected message value %q, got %q", expectedMessage, string(receivedMessages[0].Value.Must()))
	}
}
