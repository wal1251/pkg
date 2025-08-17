package serial

import (
	"encoding/json"
	"fmt"
	"io"
)

var (
	_ Encoder[any] = JSONEncode[any]
	_ Decoder[any] = JSONDecode[any]
)

func JSONEncode[T any](writer io.Writer, object T) error {
	enc := json.NewEncoder(writer)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(object); err != nil {
		return fmt.Errorf("can't serialize to JSON: %w", err)
	}

	return nil
}

func JSONDecode[T any](reader io.Reader) (T, error) {
	var object T

	if err := json.NewDecoder(reader).Decode(&object); err != nil {
		var blank T

		return blank, fmt.Errorf("can't deserialize from JSON: %w", err)
	}

	return object, nil
}
