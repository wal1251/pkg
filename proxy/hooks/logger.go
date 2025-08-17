package hooks

import (
	"context"

	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/core/presenters"
	"github.com/wal1251/pkg/proxy"
)

// LogBeforeCall возвращает хук proxy.Hook, который логирует начало вызова метода.
func LogBeforeCall(view presenters.ViewType, options presenters.ViewOptions) proxy.Hook {
	return func(ctx context.Context, object any, method string, args []any) context.Context {
		logger := logs.SubLogger(logs.FromContext(ctx),
			logs.WithElapsedTime(ctx), logs.WithMethod(method, object), logs.WithRequestID(ctx))
		logger.Info().Msg("invoking")

		if len(args) > 0 {
			logger.Trace().Msgf("args: %s", presenters.ParameterListView(args, view, options))
		}

		return logger.WithContext(ctx)
	}
}

// LogPostCall возвращает хук proxy.Hook, который логирует окончание вызова метода.
func LogPostCall(view presenters.ViewType, options presenters.ViewOptions) proxy.Hook {
	return func(ctx context.Context, _ any, _ string, results []any) context.Context {
		logger := logs.SubLogger(logs.FromContext(ctx), logs.WithElapsedTime(ctx), logs.WithRequestID(ctx))

		if proxy.HasError(results) {
			_, err := proxy.ExtractErr(results)

			logger.Err(err).Msg("returned error")
		} else {
			logger.Info().Msg("success")

			if len(results) > 0 {
				logger.Trace().Msgf("results: %s", presenters.ParameterListView(results, view, options))
			}
		}

		return ctx
	}
}

// LogPanic возвращает хук proxy.PanicHook, который перехватывает панику, логирует сообщение об ошибке и возвращает
// полученное сообщение об ошибке.
func LogPanic(view presenters.ViewType, options presenters.ViewOptions) proxy.PanicHook {
	return func(msg any, _ []byte, object any, method string, args []any) any {
		ctx := proxy.ExtractContext(args)

		logger := logs.SubLogger(logs.FromContext(ctx),
			logs.WithElapsedTime(ctx), logs.WithMethod(method, object), logs.WithRequestID(ctx))
		logger.Error().Stack().Msg("panic while performing operation")

		if len(args) > 0 {
			logger.Trace().Msgf("args: %s", presenters.ParameterListView(args, view, options))
		}

		return msg
	}
}
