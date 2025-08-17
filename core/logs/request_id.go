package logs

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

// RequestID возвращает идентификатор запроса из контекста.
func RequestID(ctx context.Context) string {
	requestID := middleware.GetReqID(ctx)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	return requestID
}

// SetRequestID возвращает новый контекст с установленным идентификатором запроса.
func SetRequestID(ctx context.Context, requestID string) context.Context {
	if requestID != "" {
		return context.WithValue(ctx, middleware.RequestIDKey, requestID)
	}

	return ctx
}
