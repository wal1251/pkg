// Package sentry
// Данный пакет позволяет инициализировать sentry клиент в приложение
package sentry

import (
	"context"
	"fmt"
	"time"

	sentry "github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"

	"github.com/wal1251/pkg/core/logs"
)

// NewSentry конструктор клиента sentry.
func NewSentry(ctx context.Context, cfg *Config, tags map[string]string) (*sentryhttp.Handler, error) {
	logger := logs.FromContext(ctx)

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.DSN,
		Debug:            cfg.Debug,
		AttachStacktrace: true,
		Environment:      cfg.Environment,
		TracesSampleRate: tracesSampleRate,
		IgnoreErrors:     []string{"context canceled", "404 Not Found"},
	})
	if err != nil {
		return nil, fmt.Errorf("can't init new centry client: %w", err)
	}

	defer sentry.Flush(time.Second)

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		for tagName, value := range tags {
			scope.SetTag(tagName, value)
		}

		scope.SetLevel(sentry.LevelError)
	})

	sentryHandler := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
		Timeout: cfg.Timeout,
	})

	logger.Info().Msg("sentry initialized")

	return sentryHandler, nil
}

func NewGRPCSentry(ctx context.Context, cfg *Config, tags map[string]string) error {
	logger := logs.FromContext(ctx)

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.DSN,
		Debug:            cfg.Debug,
		AttachStacktrace: true,
		Environment:      cfg.Environment,
		TracesSampleRate: tracesSampleRate,
		IgnoreErrors:     []string{"context canceled", "404 Not Found"},
	})
	if err != nil {
		return fmt.Errorf("can't init new sentry client: %w", err)
	}

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		for tagName, value := range tags {
			scope.SetTag(tagName, value)
		}
		scope.SetLevel(sentry.LevelError)
	})

	logger.Info().Msg("sentry initialized")

	return nil
}
