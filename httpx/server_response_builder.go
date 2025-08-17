package httpx

import (
	"context"
	"io"
	"net/http"

	"github.com/wal1251/pkg/core"
	"github.com/wal1251/pkg/core/errs"
	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/tools/serial"

	"google.golang.org/genproto/googleapis/rpc/errdetails" //nolint:typecheck
	"google.golang.org/grpc/status"
)

type ServerResponseBuilder[T any] struct {
	status  int
	value   func() T
	Header  http.Header
	encoder serial.Encoder[T]
}

func (r *ServerResponseBuilder[T]) WithContentType(contentType string) *ServerResponseBuilder[T] {
	r.Header.Set(HeaderContentType, contentType)

	return r
}

func (r *ServerResponseBuilder[T]) WithContentTypeJSON() *ServerResponseBuilder[T] {
	return r.WithContentType("application/json;charset=UTF-8").WithEncoder(serial.JSONEncode[T])
}

func (r *ServerResponseBuilder[T]) WithContentTypeXML() *ServerResponseBuilder[T] {
	return r.WithContentType("application/xml;charset=UTF-8").WithEncoder(serial.XMLEncode[T])
}

func (r *ServerResponseBuilder[T]) WithEncoder(e serial.Encoder[T]) *ServerResponseBuilder[T] {
	r.encoder = e

	return r
}

type GRPCErrorResponse struct {
	Error string `json:"error"`
}

func (r *ServerResponseBuilder[T]) WithGRPCErrorHandler() *ServerResponseBuilder[T] {
	r.WithContentTypeJSON()

	originalEncoder := r.encoder
	r.encoder = func(w io.Writer, v T) error {
		// If v is error type, handle it as GRPC error
		if err, ok := any(v).(error); ok {
			status, ok := status.FromError(err)
			if !ok {
				response := GRPCErrorResponse{Error: "0"}

				return serial.JSONEncode(w, response)
			}

			for _, detail := range status.Details() {
				if errorInfo, ok := detail.(*errdetails.ErrorInfo); ok {
					if errNum, exists := errorInfo.GetMetadata()["error_num"]; exists {
						response := GRPCErrorResponse{Error: errNum}

						return serial.JSONEncode(w, response)
					}
				}
			}

			response := GRPCErrorResponse{Error: "0"}

			return serial.JSONEncode(w, response)
		}

		// If v is not an error, use the original encoder
		return originalEncoder(w, v)
	}

	return r
}

func (r *ServerResponseBuilder[T]) WithStatus(status int) *ServerResponseBuilder[T] {
	r.status = status

	return r
}

func (r *ServerResponseBuilder[T]) WithValue(v T) *ServerResponseBuilder[T] {
	r.value = func() T { return v }

	return r
}

func (r *ServerResponseBuilder[T]) WithValueFn(f func() T) *ServerResponseBuilder[T] {
	r.value = f

	return r
}

func (r *ServerResponseBuilder[T]) Send(response http.ResponseWriter) error {
	for key, values := range r.Header {
		for _, value := range values {
			response.Header().Add(key, value)
		}
	}

	response.WriteHeader(r.status)

	if r.value == nil {
		return nil
	}

	return r.encoder.Encode(response, r.value())
}

func NewServerResponse[T any](opts ...func(*ServerResponseBuilder[T])) *ServerResponseBuilder[T] {
	builder := &ServerResponseBuilder[T]{
		Header:  make(http.Header),
		status:  ResponseStatusDefault,
		encoder: serial.VoidEncode[T],
	}

	for _, opt := range opts {
		opt(builder)
	}

	return builder
}

func ServerErrorResponses[T any](
	errToResponse core.Map[error, T],
	errToStatus *ErrorToStatusMapper,
) func(err error) *ServerResponseBuilder[T] {
	return func(err error) *ServerResponseBuilder[T] {
		return NewServerResponse[T]().
			WithContentTypeJSON().
			WithStatus(errToStatus.Status(errs.AsReason(err).Type)).
			WithValue(errToResponse.Map(err))
	}
}

func SendResponse[T any](ctx context.Context, w http.ResponseWriter, r *ServerResponseBuilder[T]) bool {
	if err := r.Send(w); err != nil {
		logs.FromContext(ctx).Err(err).Msg("failed to send http response to client")

		return false
	}

	return true
}
