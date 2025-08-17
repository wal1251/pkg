package serial

import (
	"bytes"
	"io"
)

type (
	Encoder[T any] func(w io.Writer, t T) error
	Decoder[T any] func(r io.Reader) (T, error)
)

func (e Encoder[T]) Encode(w io.Writer, t T) error {
	if e == nil {
		return nil
	}

	return e(w, t)
}

func (d Decoder[T]) Decode(reader io.Reader) (T, error) {
	var t T
	if d == nil {
		return t, nil
	}

	return d(reader)
}

func ToBytes[T any](object T, enc Encoder[T]) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := enc.Encode(buf, object); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func FromBytes[T any](raw []byte, enc Decoder[T]) (T, error) {
	return enc.Decode(bytes.NewReader(raw))
}
