package mw

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/wal1251/pkg/httpx"
)

// RequestID улучшает функциональность мидлвари middleware.RequestID - берет requestID из контекста запроса и помещает
// его в заголовок ответа.
func RequestID() httpx.Middleware {
	return httpx.MiddlewareFn(func(response http.ResponseWriter, request *http.Request, next http.Handler) {
		requestID := request.Context().Value(middleware.RequestIDKey)
		if requestID != nil {
			if id, ok := requestID.(string); ok {
				response.Header().Set(middleware.RequestIDHeader, id)
			}
		}
		next.ServeHTTP(response, request)
	}).Middleware()
}

// AnnotateRequestContext добавляет заголовки запроса в metadata контекста для grpc запроса.
func AnnotateRequestContext() httpx.Middleware {
	return httpx.MiddlewareFn(func(response http.ResponseWriter, request *http.Request, next http.Handler) {
		next.ServeHTTP(response, request.WithContext(httpx.AnnotateContext(request)))
	}).Middleware()
}
