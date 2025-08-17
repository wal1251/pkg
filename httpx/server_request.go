package httpx

import (
	"bytes"
	"net/http"

	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/tools/serial"
)

type (
	ServerRequest[T any] struct {
		*http.Request
		decoder serial.Decoder[T]
		Value   T
	}
)

func (r *ServerRequest[T]) WithAutoDecoder() *ServerRequest[T] {
	if Header(r.Header).HasJSONContent() {
		return r.WithDecoder(serial.JSONDecode[T])
	}
	if Header(r.Header).HasXMLContent() {
		return r.WithDecoder(serial.XMLDecode[T])
	}

	return r.WithDecoder(serial.VoidDecode[T])
}

func (r *ServerRequest[T]) WithDecoder(d serial.Decoder[T]) *ServerRequest[T] {
	r.decoder = d

	return r
}

func (r *ServerRequest[T]) Decode() error {
	defer func() {
		if r.Body != nil {
			if err := r.Body.Close(); err != nil {
				logs.FromContext(r.Context()).Err(err).Msg("failed to close request body")
			}
		}
	}()

	raw, err := (*Request)(r.Request).ReadBody()
	if err != nil {
		return err
	}

	request, err := r.decoder.Decode(bytes.NewBuffer(raw))
	if err != nil {
		return err
	}

	r.Value = request

	return nil
}

func (r *ServerRequest[T]) ReadFile(name string) (*MultipartFile, error) {
	return NewMultipartFile(r.Request, name)
}

func NewServerRequest[T any](r *http.Request) *ServerRequest[T] {
	return (&ServerRequest[T]{Request: r}).WithAutoDecoder()
}
