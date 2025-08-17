// Package otelx.
// Данный пакет предназначен для создания провайдера трассировки на основе OpenTelemetry.
package otelx

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

func NewOTELTraceProvider(ctx context.Context, cfg *Config) (trace.TracerProvider, error) {
	// Создаем экспортер для отправки данных трейсов в формате OTLP.
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(cfg.ExporterEndpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Создаем ресурс, который будет содержать информацию о сервисе.
	resource, err := sdkresource.New(ctx,
		sdkresource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.DeploymentEnvironmentKey.String(cfg.Env),
		),
		sdkresource.WithSchemaURL(semconv.SchemaURL),
		sdkresource.WithContainer(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP resource: %w", err)
	}

	// Создаем сэмплер для отправки определенного процента трейсов.
	// Если родительский спан начал отслеживаться, то будут отслеживаться и все дочерние спаны.
	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.SamplerTraceRatio))

	// Создаем провайдер трассировки.
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(sampler),
	)

	return tracerProvider, nil
}
