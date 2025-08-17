package serial

import (
	"encoding/xml"
	"fmt"
	"io"
)

var (
	_ Encoder[any] = XMLEncode[any]
	_ Decoder[any] = XMLDecode[any]
)

func XMLEncode[T any](writer io.Writer, object T) error {
	enc := xml.NewEncoder(writer)

	if err := enc.Encode(object); err != nil {
		return fmt.Errorf("can't serialize to XML: %w", err)
	}

	return nil
}

func XMLDecode[T any](reader io.Reader) (T, error) {
	var object T

	if err := xml.NewDecoder(reader).Decode(&object); err != nil {
		var blank T

		return blank, fmt.Errorf("can't deserialize from XML: %w", err)
	}

	return object, nil
}
