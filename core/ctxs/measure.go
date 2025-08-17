package ctxs

import (
	"context"
	"time"
)

const (
	// StartTimeMeasureContextKey ключ хранения начала замера длительности операции. См. StartMeasureContext и ElapsedFromContext.
	StartTimeMeasureContextKey = "CLOCK-StartTimeMeasure"
)

// StartMeasureContext начинает замер длительности операции, выполнить замер можно с помощью GetElapsedTime.
func StartMeasureContext(ctx context.Context) context.Context {
	return ValuePut(ctx, StartTimeMeasureContextKey, time.Now().UnixMicro())
}

// ElapsedFromContext получает из контекста прошедшее время с момента начала замера с помощью StartTimeMeasure.
func ElapsedFromContext(ctx context.Context) time.Duration {
	return time.Since(time.UnixMicro(ValueGet[int64](ctx, StartTimeMeasureContextKey)))
}
