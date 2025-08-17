package grpcx

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// AddValueToGrpcCtx добавляет указанную пару ключ-значение в метаданные gRPC-контекста.
// Функция создает новый исходящий контекст с метаданными (metadata.MD),
// который используется для передачи данных между клиентом и сервером через gRPC.
func AddValueToGrpcCtx(ctx context.Context, key, value string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, key, value)
}
