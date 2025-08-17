package images

import (
	"io"
)

type (
	Preview struct {
		Name    string
		Format  string
		Content io.ReadSeeker
		Size    int64
	}

	PreviewOptions struct {
		JpegQuality int
	}
)
