package logs

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/wal1251/pkg/core/ctxs"
	"github.com/wal1251/pkg/core/presenters"
)

const (
	InstanceIDTag      Tag = "instance_id"       // Идентифицирует экземпляр приложения.
	MethodTag          Tag = "method"            // Вызванный метод.
	ElapsedTag         Tag = "elapsed"           // Время, прошедшее с начала получения запроса.
	RequestIDTag       Tag = "request_id"        // Уникальный идентификатор входящего запроса.
	FromTag            Tag = "from"              // Адрес вызывающего.
	UserIDTag          Tag = "user_id"           // Идентификатор авторизованного пользователя.
	UserAuthorityTag   Tag = "user_authority"    // Полномочия авторизованного пользователя.
	TokenIDTag         Tag = "token_id"          // Идентификатор токена авторизованного пользователя.
	ComponentTag       Tag = "component"         // Идентифицирует вызываемый компонент.
	ComponentLabelTag  Tag = "component_label"   // Идентифицирует экземпляр вызываемого компонента.
	ComponentCallIDTag Tag = "component_call_id" // Идентифицирует конкретный запрос к компоненту.
	PathTag            Tag = "path"              // URL запроса.
	HeadersTag         Tag = "headers"           // Заголовки запроса.
)

// Tag структурный тэг логов.
type Tag string

func TagValueArray[T any](s []T, f func(*zerolog.Array, T)) *zerolog.Array {
	arr := zerolog.Arr()

	for _, t := range s {
		f(arr, t)
	}

	return arr
}

func TagStringArray(s []string) *zerolog.Array {
	return TagValueArray(s, func(a *zerolog.Array, t string) { a.Str(t) })
}

// Option возвращает функциональную опцию логера, которая применит к контексту логера тэг t с указанным значением value.
func (t Tag) Option(value any) LoggerOption {
	return func(logCtx zerolog.Context) zerolog.Context {
		if value == nil {
			return logCtx
		}

		if v, ok := value.(presenters.StringViewer); ok {
			return logCtx.Str(string(t), v.StringView(presenters.ViewLogs, presenters.ViewOptions{}))
		}

		switch typedValue := value.(type) {
		case time.Duration:
			return logCtx.Dur(string(t), typedValue)
		case string:
			return logCtx.Str(string(t), typedValue)
		case fmt.Stringer:
			return logCtx.Stringer(string(t), typedValue)
		case *zerolog.Array:
			return logCtx.Array(string(t), typedValue)
		default:
			return logCtx.Interface(string(t), value)
		}
	}
}

// Value возвращает функциональную опцию события логера, которая применит к событию тэг t с указанным значением value.
func (t Tag) Value(value any) EventOption {
	return func(event *zerolog.Event) *zerolog.Event {
		if value == nil {
			return event
		}

		if v, ok := value.(presenters.StringViewer); ok {
			return event.Str(string(t), v.StringView(presenters.ViewLogs, presenters.ViewOptions{}))
		}

		switch typedValue := value.(type) {
		case time.Duration:
			return event.Dur(string(t), typedValue)
		case string:
			return event.Str(string(t), typedValue)
		case fmt.Stringer:
			return event.Stringer(string(t), typedValue)
		case *zerolog.Array:
			return event.Array(string(t), typedValue)
		default:
			return event.Interface(string(t), value)
		}
	}
}

// EventWith применяет к событию логера набор указанных функциональных опций.
func EventWith(event *zerolog.Event, modifiers ...EventOption) *zerolog.Event {
	for _, modifier := range modifiers {
		event = modifier.ApplyTo(event)
	}

	return event
}

// EventComponentCall добавляет к событию логера набор тегов, присущих для вызова компоненты.
func EventComponentCall(name, label string) EventOption {
	callID := uuid.NewString()

	return func(event *zerolog.Event) *zerolog.Event {
		return EventWith(event,
			ComponentTag.Value(name).If(name != ""),
			ComponentLabelTag.Value(label).If(name != "" && label != ""),
			ComponentCallIDTag.Value(callID).If(name != ""),
		)
	}
}

// WithElapsedTime возвращает функциональную опцию логера, которая установит тэг логера с длительностью выполнения
// операции, см. ElapsedTag. Замер времени запускается функцией ctxs.StartMeasureContext().
func WithElapsedTime(ctx context.Context) LoggerOption {
	elapsed := ctxs.ElapsedFromContext(ctx)

	return func(z zerolog.Context) zerolog.Context {
		if elapsed == 0 {
			return z
		}

		return ElapsedTag.Option(elapsed).ApplyTo(z)
	}
}

// WithRequestID возвращает функциональную опцию логера, которая установит тег логгера с уникальным идентификатором
// запроса, см. RequestIDTag. Идентификатор будет прочитан из указанного контекста ctx.
func WithRequestID(ctx context.Context) LoggerOption {
	return RequestIDTag.Option(RequestID(ctx))
}

// WithMethod возвращает функциональную опцию логера, которая применит и установит тег с названием метода,
// сформированный из значений method и object (в формате {object}::{method}).
func WithMethod(method string, object any) LoggerOption {
	return func(z zerolog.Context) zerolog.Context {
		if object == nil {
			return MethodTag.Option(method).ApplyTo(z)
		}

		return MethodTag.Option(fmt.Sprintf("%T::%s", object, method)).ApplyTo(z)
	}
}
