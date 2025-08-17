package images

import (
	"fmt"
	"io"

	"github.com/disintegration/imaging"
)

var _ Transformer = (*Encode)(nil)

type (
	Encode struct {
		JpegQuality int
		Destination io.Writer
		Format      string
	}
)

func (t Encode) Perform(img Image) (Image, error) {
	ext := t.Format
	if ext == "" {
		ext = "png"
	}

	format, err := imaging.FormatFromExtension(ext)
	if err != nil {
		return Image{}, fmt.Errorf("can't autodetect image format from extension: %w", err)
	}

	options := make([]imaging.EncodeOption, 0)
	if t.JpegQuality != 0 {
		options = append(options, imaging.JPEGQuality(t.JpegQuality))
	}

	if err = imaging.Encode(t.Destination, img, format, options...); err != nil {
		return Image{}, fmt.Errorf("can't encode imege to format %s: %w", format, err)
	}

	return img, nil
}
