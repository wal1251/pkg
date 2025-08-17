package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/wal1251/pkg/tools/collections"
	"github.com/wal1251/pkg/tools/serial"
)

const (
	DefaultPartitions        = 16
	DefaultReplicationFactor = 1
)

type (
	ValueProvider func() ([]byte, error)

	Message struct {
		Topic     string
		Partition *int32
		Key       []byte
		Value     ValueProvider
		Headers   collections.MultiMap[string, []byte]
		OnAck     func(context.Context)
	}

	TopicMetadata struct {
		Name              string
		Partitions        int
		ReplicationFactor int
	}
)

func (p ValueProvider) Get() ([]byte, error) {
	if p != nil {
		return p()
	}

	return nil, nil
}

func (p ValueProvider) Must() []byte {
	v, err := p.Get()
	if err != nil {
		panic(err)
	}

	return v
}

func (m *Message) WithAck(ack func(context.Context)) *Message {
	m.OnAck = ack

	return m
}

func (m *Message) WithValue(value ValueProvider) *Message {
	m.Value = value

	return m
}

func (m *Message) WithPartition(partition int32) *Message {
	m.Partition = &partition

	return m
}

func (m *Message) WithKeyBytes(key []byte) *Message {
	if len(key) != 0 {
		m.Key = key
	}

	return m
}

func (m *Message) WithKey(key any) *Message {
	return m.WithKeyBytes([]byte(fmt.Sprint(key)))
}

func (m *Message) WithHeaderBytes(key string, value []byte) *Message {
	if len(value) == 0 {
		return m
	}

	if m.Headers == nil {
		m.Headers = make(collections.MultiMap[string, []byte])
	}

	m.Headers.Append(key, value)

	return m
}

func (m *Message) WithHeader(key string, value string) *Message {
	return m.WithHeaderBytes(key, []byte(value))
}

func (m *Message) Header(key string) string {
	return string(m.HeaderBytes(key))
}

func (m *Message) HeaderBytes(key string) []byte {
	if len(m.Headers) == 0 {
		return nil
	}

	value, ok := m.Headers[key]
	if !ok {
		return nil
	}

	if len(value) == 0 {
		return nil
	}

	return value[0]
}

func (m *Message) Ack(ctx context.Context) {
	if m.OnAck == nil {
		return
	}

	m.OnAck(ctx)
}

func ValueJSON(value any) ValueProvider {
	return func() ([]byte, error) { return serial.ToBytes(value, serial.JSONEncode[any]) }
}

func ValueBytes(value []byte) ValueProvider {
	return func() ([]byte, error) { return value, nil }
}

func ValueJSONToTyped[T any](val ValueProvider) (T, error) {
	raw, err := val.Get()
	if err != nil {
		var blank T

		return blank, err
	}

	return serial.FromBytes(raw, serial.JSONDecode[T])
}

func Milliseconds(duration time.Duration) int {
	return int(duration.Milliseconds())
}

func MakeTopicMetadata(name string) TopicMetadata {
	return TopicMetadata{
		Name:              name,
		Partitions:        DefaultPartitions,
		ReplicationFactor: DefaultReplicationFactor,
	}
}
