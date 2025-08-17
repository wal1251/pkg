// Package logs предоставляет функции и вспомогательные примитивы для работы с логами в приложении.
//
// Пакет является фасадом над библиотекой zerolog.
//
// Не приветствуется явная передача логера через аргумент, логер должен храниться в контексте приложения, при
// необходимости извлекаться из него.
package logs

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/wal1251/pkg/core"
)

type (
	// LoggerOption функциональная опция логера.
	LoggerOption func(zerolog.Context) zerolog.Context

	// EventOption функциональная опция события логера.
	EventOption func(*zerolog.Event) *zerolog.Event
)

// ApplyTo применяет функциональную опцию к указанному контексту логера и возвращает обновленный контекст логера.
func (o LoggerOption) ApplyTo(z zerolog.Context) zerolog.Context {
	if o == nil {
		return z
	}

	return o(z)
}

// ApplyTo применяет функциональную опцию к указанному событию логера и возвращает указатель на это событие.
func (m EventOption) ApplyTo(e *zerolog.Event) *zerolog.Event {
	if m == nil {
		return e
	}

	return m(e)
}

// If вернет новую функциональную опцию логера, которая применяется к логеру только если condition равно true.
func (m EventOption) If(condition bool) EventOption {
	return func(event *zerolog.Event) *zerolog.Event {
		if condition {
			m.ApplyTo(event)
		}

		return event
	}
}

// To создает новое событие логера из предоставленной функции и применяет себя к новому событию, возвращает указатель на
// это событие. Например, можно использовать так:
//
//	logger := FromContext(ctx)
//	var logContext EventOption = getLocalLoggerContext()
//	log := logContext.To(logger.Debug).Msg("connection established")
//
// Может быть полезным, чтобы формировать локальный (в рамках вызова функции) контекст логера.
func (m EventOption) To(newEvent func() *zerolog.Event) *zerolog.Event {
	return m.ApplyTo(newEvent())
}

// Logger возвращает новый экземпляр логера с указанными опциями LoggerOption.
func Logger(cfg *Config, options ...LoggerOption) zerolog.Logger {
	var output io.Writer = os.Stdout
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{Out: output, TimeFormat: time.RFC3339}
	}

	lvl, _ := zerolog.ParseLevel(cfg.Level)

	return Options(options...).ApplyTo(zerolog.New(output).Level(lvl).With().Timestamp()).Logger()
}

// SubLogger возвращает новый сублогер, наследованный от указанного, с примененными функциональными опциями.
func SubLogger(logger *zerolog.Logger, options ...LoggerOption) zerolog.Logger {
	return Options(options...).ApplyTo(logger.With()).Logger()
}

// ToContext помещает указанный указатель на логер в контекст.
func ToContext(ctx context.Context, logger *zerolog.Logger) context.Context {
	return logger.WithContext(ctx)
}

// FromContext получает ранее сохраненный в контексте указатель на логер. Если ранее не был установлен, вернет дефолтный
// логер. См. zerolog.Ctx.
func FromContext(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

// Options получает новую опцию логера, которая включает в себя все указанные опции в качестве аргументов.
func Options(options ...LoggerOption) LoggerOption {
	return func(z zerolog.Context) zerolog.Context {
		for _, option := range options {
			z = option(z)
		}

		return z
	}
}

// LocalContext получает локальный контекст вызова компоненты для модификации сообщения логера.
// Например:
//
//	var component core.Component = self() // сам объект компонента.
//	logs.LocalContext(component).To(logs.FromContext(ctx).Debug).Msg("connection established")
//
// Полезно для логирования вызова методов компонента.
func LocalContext(component core.Component) EventOption {
	if component == nil {
		return func(event *zerolog.Event) *zerolog.Event {
			return event
		}
	}

	return EventComponentCall(component.Name(), component.Label())
}
