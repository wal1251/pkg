package hooks

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/wal1251/pkg/proxy"
)

const (
	UseCasesLayer = "UseCase"
	StorageLayer  = "Storage"
)

// StartSpanBeforeCall возвращает хук proxy.Hook, который инициирует Span перед вызовом метода.
func StartSpanBeforeCall(layer string) proxy.Hook {
	return func(ctx context.Context, _ any, method string, _ []any) context.Context {
		tracer := otel.Tracer("rmr-pkg")

		spanName := fmt.Sprintf("%s: %s", layer, method)
		ctx, _ = tracer.Start(ctx, spanName) // nolint:spancheck

		return ctx
	}
}

// EndSpanPostCall возвращает хук proxy.Hook, который завершает Span после вызова метода.
func EndSpanPostCall(ctx context.Context, _ any, _ string, results []any) context.Context {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	if proxy.HasError(results) {
		_, err := proxy.ExtractErr(results)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "success")
	}

	return ctx
}

// EndSpanPanic возвращает хук proxy.PanicHook, который перехватывает панику во время выполнения метода,
// хук завершает Span и записывает информацию о панике.
func EndSpanPanic(msg any, _ []byte, _ any, _ string, args []any) any {
	ctx := proxy.ExtractContext(args)

	span := trace.SpanFromContext(ctx)
	defer span.End()

	spanMsg := fmt.Sprintf("panic: %s", msg)
	span.SetStatus(codes.Error, spanMsg)

	return msg
}
