// Package events позволяет эмитировать и обрабатывать события, подключаясь к ent hooks.
package events

import (
	"context"
	"fmt"
	"strings"

	"entgo.io/ent"
	"github.com/google/uuid"

	"github.com/wal1251/pkg/core/bus"
	"github.com/wal1251/pkg/db/entx"
)

const (
	OperationCreate          Operation = 1 << iota // Запись создана.
	OperationUpdate                                // Запись обновлена.
	OperationDelete                                // Запись удалена.
	OperationDeleteFlagSet                         // Установлен флаг удаления.
	OperationDeleteFlagUnset                       // Сброшен флаг удаления.

	OperationCreateName          = "CREATE"
	OperationUpdateName          = "UPDATE"
	OperationDeleteName          = "DELETE"
	OperationDeleteFlagSetName   = "DELETE_fLAG_SET"
	OperationDeleteFlagUnsetName = "DELETE_fLAG_UNSET"
)

type (
	Operation uint // Операция, произведенная над записью.

	// Event событие, сигнализирующее об изменении модели.
	Event struct {
		Op         Operation      // Выполненная операция (создание, модификация, удаление и т.д.).
		ID         uuid.UUID      // ID сущности (если нет, тогда uuid.Nil).
		Type       string         // Тип модели, с которым связано событие.
		Attributes map[string]any // Изменяемые атрибуты модели.
	}

	// Publisher содержит список подписчиков, которые будут получать события.
	//
	// Пример использования:
	//
	// var db *ent.Client = // получаем клиент ent...
	// publisher := events.NewPublisher()
	// db.Use(events.NewEventHook(publisher.DbConsumer()))
	//
	// publisher.Subscribe(context.Background(), AfterUserCreate) // AfterUserCreate должен удовлетворять интерфейсу core.Subscriber.
	//
	Publisher struct {
		subscribers []bus.Subscriber[Event]
	}
)

func (e Event) String() string {
	attributes := make([]string, 0, len(e.Attributes))
	for k, v := range e.Attributes {
		attributes = append(attributes, fmt.Sprintf("%s=%v", k, v))
	}
	joinedAttributes := strings.Join(attributes, ", ")

	if e.ID == uuid.Nil {
		return fmt.Sprintf("%s: %s. Changed attributes: %s", e.Type, e.Op, joinedAttributes)
	}

	return fmt.Sprintf("%s(%v): %s. Changed attributes: %s", e.Type, e.ID, e.Op, joinedAttributes)
}

func (p *Publisher) Subscribe(_ context.Context, subscriber bus.Subscriber[Event]) error {
	p.subscribers = append(p.subscribers, subscriber)

	return nil
}

// DBConsumer выполняет роль подписчика-родителя, который регистрируется в
// качестве хука мутации.
// Он принимает события и передает их зарегистрированным подписчикам Publisher'а.
//
// Пример использования см. Publisher.
func (p *Publisher) DBConsumer() bus.SubscriberFn[Event] {
	return func(ctx context.Context, events ...Event) error {
		for _, subscriber := range p.subscribers {
			if err := subscriber.Publish(ctx, events...); err != nil {
				return err
			}
		}

		return nil
	}
}

func NewPublisher(subscribers ...bus.Subscriber[Event]) *Publisher {
	return &Publisher{
		subscribers: subscribers,
	}
}

// NewDBEventHook публикует события ent hooks. Адаптер для core.Subscriber.
//
// Например:
//
//	var db *ent.Client = // получаем клиент ent...
//	db.Use(Publish(core.SubscriberFn[Event](func(ctx context.Context, events ...Event) error {
//		for _, event := range events {
//			// обрабатываем событие...
//		}
//
//		return nil
//	})))
//
// .
func NewDBEventHook(subscriber bus.Subscriber[Event]) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
			value, err := next.Mutate(ctx, mutation)
			if err != nil {
				return nil, fmt.Errorf("failed to mutate object: %w", err)
			}

			if err = subscriber.Publish(ctx, MakeEvent(mutation)); err != nil {
				return nil, err
			}

			return value, nil
		})
	}
}

func NewDBEventHooks(subscribers ...bus.SubscriberFn[Event]) []ent.Hook {
	eventHooks := make([]ent.Hook, 0, len(subscribers))
	for _, subscriber := range subscribers {
		eventHooks = append(eventHooks, NewDBEventHook(subscriber))
	}

	return eventHooks
}

// MakeEvent конструктор Event из ent.Mutation.
func MakeEvent(mutation ent.Mutation) Event {
	var id uuid.UUID

	attributes := make(map[string]any)

	for _, fieldName := range mutation.Fields() {
		if fieldValue, ok := mutation.Field(fieldName); ok {
			attributes[fieldName] = fieldValue
		}
	}

	if v, ok := mutation.(interface{ ID() (uuid.UUID, bool) }); ok {
		if mID, exists := v.ID(); exists {
			id = mID
		}
	}

	operation := MakeOperation(mutation.Op())
	if attributes[entx.FieldIsDeleted] == true {
		operation = operation.With(OperationDeleteFlagSet)
	}

	if !mutation.Op().Is(ent.OpCreate) && attributes[entx.FieldIsDeleted] == false {
		operation = operation.With(OperationDeleteFlagUnset)
	}

	return Event{
		ID:         id,
		Op:         operation,
		Type:       mutation.Type(),
		Attributes: attributes,
	}
}

func (o Operation) String() string {
	operations := map[Operation]string{
		OperationCreate:          OperationCreateName,
		OperationUpdate:          OperationUpdateName,
		OperationDelete:          OperationDeleteName,
		OperationDeleteFlagSet:   OperationDeleteFlagSetName,
		OperationDeleteFlagUnset: OperationDeleteFlagUnsetName,
	}

	names := make([]string, 0, 1)

	for operation, name := range operations {
		if o.Is(operation) {
			names = append(names, name)
		}
	}

	return strings.Join(names, "&")
}

func (o Operation) Any(operations ...Operation) bool {
	for _, operation := range operations {
		if o&operation != 0 {
			return true
		}
	}

	return false
}

func (o Operation) Is(target Operation) bool {
	return o&target != 0
}

func (o Operation) With(target Operation) Operation {
	return o | target
}

func (o Operation) Without(target Operation) Operation {
	return o & ^target
}

func MakeOperation(entOp ent.Op) Operation {
	switch {
	case entOp.Is(ent.OpCreate):
		return OperationCreate
	case entOp.Is(ent.OpUpdate), entOp.Is(ent.OpUpdateOne):
		return OperationUpdate
	case entOp.Is(ent.OpDelete), entOp.Is(ent.OpDeleteOne):
		return OperationDelete
	}

	return Operation(0)
}
