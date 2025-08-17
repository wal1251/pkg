package api

import (
	"io"

	"github.com/elastic/go-elasticsearch/v8/esutil"
)

var _ io.Reader = (*ErrReader)(nil)

type ErrReader struct {
	Err error
}

func (r *ErrReader) Read([]byte) (int, error) {
	return 0, r.Err
}

func NewErrReader(err error) *ErrReader {
	return &ErrReader{err}
}

func NewJSONReader(value any) io.Reader {
	if value == nil {
		return nil
	}

	return esutil.NewJSONReader(value)
}
