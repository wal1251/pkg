package serial

import (
	"fmt"
	"io"
)

var (
	_ Encoder[[]byte] = BytesWrite
	_ Decoder[[]byte] = BytesRead
)

func BytesWrite(w io.Writer, t []byte) error {
	_, err := w.Write(t)

	return fmt.Errorf("can't write bytes: %w", err)
}

func BytesRead(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("can't read bytes: %w", err)
	}

	return data, nil
}
