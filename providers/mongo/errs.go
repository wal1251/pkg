package mongo

import (
	"errors"
)

var (
	ErrIDNotFoundInDocument = errors.New("id not found in document")
	ErrInvalidIDInDocument  = errors.New("field '_id' is not a valid ObjectID")

	ErrStringValueConversionFailed = errors.New("value to string conversion failed")

	ErrDocumentValueNotFound = errors.New("document value not found")
)
