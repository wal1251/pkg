package events

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/wal1251/pkg/providers/kafka"
)

type TestEvent struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestEventRead(t *testing.T) {
	// create a test context
	ctx := context.TODO()

	// create a test message
	testEvent := TestEvent{
		Name: "John Doe",
		Age:  30,
	}
	testEventBytes, _ := json.Marshal(testEvent)
	testMessage := kafka.Message{
		Value: kafka.ValueJSON(testEventBytes),
	}

	// define the accept function
	accept := func(ctx context.Context, event TestEvent) {
		if event.Name != testEvent.Name {
			t.Errorf("expected name %s, but got %s", testEvent.Name, event.Name)
		}
		if event.Age != testEvent.Age {
			t.Errorf("expected age %d, but got %d", testEvent.Age, event.Age)
		}
	}

	// call the EventRead function with the test message and accept function
	reader := EventRead[TestEvent](accept)
	reader(ctx, &testMessage)
}
