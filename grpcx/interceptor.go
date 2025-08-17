package grpcx

import (
	"context"
	"fmt"
	"strings"

	"github.com/wal1251/pkg/tools/acceptlanguage"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/wal1251/pkg/core/errs"
)

// UserInfoKey - ключ для хранения информации о пользователе в контексте.
const UserInfoKey = "userInfo"

// UserInfoClientInterceptor добавляет информацию о пользователе в метаданные gRPC-запроса.
func UserInfoClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		clientConn *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Извлекаем информацию о пользователе из контекста
		userInfo := ctx.Value(UserInfoKey)
		if userInfo != nil {
			claims, ok := userInfo.(map[string]interface{})
			if ok {
				md := metadata.New(nil)
				for key, value := range claims {
					md.Append(key, fmt.Sprintf("%v", value))
				}
				ctx = metadata.NewOutgoingContext(ctx, md)
			}
		}
		// Вызываем RPC
		return invoker(ctx, method, req, reply, clientConn, opts...)
	}
}

// UserInfoServerInterceptor извлекает информацию о пользователе из метаданных и добавляет ее в контекст.
func UserInfoServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		neededKeys := []string{"user_id", "phone_number", "name", "session_id"}
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			userInfo := make(map[string]interface{})
			for key, values := range md {
				for _, neededKeys := range neededKeys {
					if key == neededKeys {
						userInfo[key] = values[0]
					}
				}
			}
			ctx = context.WithValue(ctx, UserInfoKey, userInfo) //nolint:revive,staticcheck //FIXME
		}

		return handler(ctx, req)
	}
}

// SecureErrorInterceptor Интерцептор для обработки ошибок GRPC с цифровым кодом и возвращением ошибки в формате grpc.Status.
func SecureErrorInterceptor(serviceID int) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		h, err := handler(ctx, req)
		if err != nil { // nolint:nestif
			if _, ok := status.FromError(err); ok {
				err = errs.ErrSystemFailure
			}
			err = NewDetailedGrpcError(serviceID, err)

			return nil, err
		}

		return h, nil
	}
}

// ErrorInterceptor Интерсептор для обработки ошибок GRPC.
func ErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		h, err := handler(ctx, req)
		if err != nil { // nolint:nestif
			if errFromGRPC, ok := status.FromError(err); ok {
				return nil, errFromGRPC.Err() // nolint:wrapcheck
			}
			err = NewGrpcError(err)

			return nil, err
		}

		return h, nil
	}
}

// SentryInterceptor is a gRPC UnaryServerInterceptor for capturing errors and additional context into Sentry.
func SentryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	h, err := handler(ctx, req)
	if err != nil { // nolint:nestif
		// Extract gRPC status from the error
		status, _ := status.FromError(err)

		// Capture the error and additional context with Sentry
		sentry.WithScope(
			func(scope *sentry.Scope) {
				// Add gRPC method to Sentry tags
				scope.SetTag("grpc.method", info.FullMethod)

				// Add gRPC status code to Sentry tags
				scope.SetTag("grpc.status_code", status.Code().String())

				// Add metadata from the incoming context to Sentry tags
				if md, ok := metadata.FromIncomingContext(ctx); ok {
					for key, values := range md {
						// If the key is related to authorization, mask it
						if key == "authorization" || key == "token" {
							scope.SetTag(key, "REDACTED")
						} else {
							for _, value := range values {
								scope.SetTag(key, value)
							}
						}
					}
				}

				// Capture the exception with all the configured scope
				sentry.CaptureException(err)
			},
		)
	}

	return h, err
}

// RequestIDServerInterceptor извлекает x-request-id из metadata grpc запроса и добавляет в контекст.
func RequestIDServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			reqID := uuid.New().String()
			v := md.Get(strings.ToLower(middleware.RequestIDHeader))
			if len(v) > 0 {
				reqID = v[0]
			}

			ctx = context.WithValue(ctx, middleware.RequestIDKey, reqID)
		}

		return handler(ctx, req)
	}
}

// AcceptLanguageClientInterceptor добавляет информацию о языке в метаданные gRPC-запроса.
func AcceptLanguageClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		clientConn *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Извлекаем информацию о языке из контекста
		acceptLanguage := ctx.Value(acceptlanguage.AcceptLanguageKey)
		if acceptLanguage != nil {
			alSubtag, ok := acceptLanguage.(string)
			if ok {
				md := metadata.New(map[string]string{
					"accept_language": alSubtag,
				})
				ctx = metadata.NewOutgoingContext(ctx, md)
			}
		}
		// Вызываем RPC
		return invoker(ctx, method, req, reply, clientConn, opts...)
	}
}
