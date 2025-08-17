package serial

import "io"

var (
	_ Encoder[any] = VoidEncode[any]
	_ Decoder[any] = VoidDecode[any]
)

func VoidEncode[T any](io.Writer, T) error {
	return nil
}

func VoidDecode[T any](io.Reader) (T, error) {
	var blank T

	return blank, nil
}
