package mw

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"

	"github.com/wal1251/pkg/httpx"
)

// OTELContextPropagator автоматически начинает новый span для каждого входящего HTTP запроса.
func OTELContextPropagator() httpx.Middleware {
	return otelhttp.NewMiddleware("HTTP Server",
		otelhttp.WithPropagators(propagation.TraceContext{}),
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return fmt.Sprintf("%s: %s %s", operation, r.Method, r.URL)
		}),
	)
}
