package otelx_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/wal1251/pkg/httpx/mw"
	"github.com/wal1251/pkg/providers/otelx"
)

// ExampleOTELTraceProvider демонстрирует создание провайдера трассировки и создание спанов.
func ExampleOTELTraceProvider() {
	ctx := context.Background()

	// Предварительно запускаем Jaeger в контейнере:
	// docker run --name jaeger \
	//  -e COLLECTOR_OTLP_ENABLED=true \
	//  -p 16686:16686 \
	//  -p 4317:4317 \
	//  -p 4318:4318 \
	//  jaegertracing/all-in-one:latest

	// Создаем конфигурацию для провайдера трассировки.
	cfg := &otelx.Config{
		ExporterEndpoint:  "http://localhost:4318",
		ServiceName:       "exampleService",
		Env:               "development",
		SamplerTraceRatio: 1.0, // Отправка 100% трейсов
	}

	// Создание провайдера трассировки.
	traceProvider, err := otelx.NewOTELTraceProvider(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create OTEL trace provider: %v", err)
	}

	// Регистрируем провайдер трассировки глобально.
	otel.SetTracerProvider(traceProvider)

	// Получаем tracer для создания спанов.
	tracer := otel.Tracer("exampleTracer")

	// Создаем родительский спан для нашей задачи.
	ctx, span := tracer.Start(ctx, "exampleTask")
	// Завершаем спан после выполнения задачи.
	defer span.End()

	// Добавляем атрибуты к спану.
	span.SetAttributes(attribute.String("exampleKey", "exampleValue"))
	// Добавляем события (логи) к спану.
	span.AddEvent("exampleEvent")
	// Устанавливаем статус спана.
	span.SetStatus(codes.Error, "exampleError")

	// Вызываем вложенную функцию, передавая контекст.
	executeSubTask(ctx)
}

// executeSubTask демонстрирует выполнение подзадачи с трассировкой.
func executeSubTask(ctx context.Context) {
	tracer := otel.Tracer("exampleTracer")
	// Создаем дочерний спан для подзадачи.
	_, subSpan := tracer.Start(ctx, "executeSubTask")
	defer subSpan.End()

	// Имитация выполнения подзадачи.
	time.Sleep(50 * time.Millisecond)
}

// ExampleHTTPServerMiddlewareWithContextPropagation демонстрирует использование
// middleware получения данных трассировки из HTTP запроса.
func ExampleHTTPServerMiddlewareWithContextPropagation() {
	// Предполагается, что провайдер трассировки уже глобально зарегистрирован.

	// Middleware для получения данных трассировки из HTTP запроса.
	middleware := mw.OTELContextPropagator()

	http.Handle("/", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})))
}

// ExampleHTTPClientWithContextPropagation демонстрирует использование
// механизма распространения контекста трассировки при отправке HTTP запроса.
func ExampleHTTPClientWithContextPropagation() {
	// Предполагается, что провайдер трассировки уже глобально зарегистрирован.

	// Создаем клиент для отправки запросов с данными трассировки.
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport, otelhttp.WithPropagators(propagation.TraceContext{}))}

	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://localhost:8080", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Отправляем запрос.
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()
}

// ExampleGRPCServerStatsHandlerWithContextPropagation демонстрирует использование
// middleware для получения данных трассировки из GRPC запроса.
func ExampleGRPCServerStatsHandlerWithContextPropagation() {
	// Предполагается, что провайдер трассировки уже глобально зарегистрирован.

	_ = grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(
				otelgrpc.WithPropagators(propagation.TraceContext{}),
			),
		),
	)
}

// ExampleGRPCClientWithContextPropagation демонстрирует использование
// механизма распространения контекста трассировки при отправке GRPC запроса.
func ExampleGRPCClientWithContextPropagation() {
	// Предполагается, что провайдер трассировки уже глобально зарегистрирован.

	// Создаем коннект к GRPC серверу для отправки запросов с данными трассировки.
	_, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(
				otelgrpc.WithPropagators(propagation.TraceContext{}))))
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to connect: %v", err))
	}
}
