package images

import (
	"fmt"
	"image"
	"io"

	"github.com/disintegration/imaging"
)

var _ Imageable = (*Image)(nil)

type (
	Image struct {
		image.Image
		Metadata map[string]any
	}
)

func MakeImage() Image {
	return Image{
		Metadata: make(map[string]any),
	}
}

func (i *Image) Decode(input io.Reader) error {
	img, err := imaging.Decode(input)
	if err != nil {
		return fmt.Errorf("can't devcode image: %w", err)
	}

	i.Image = img

	return nil
}

func (i *Image) Transform(t func(image.Image) image.Image) Image {
	return Image{
		Image:    t(i),
		Metadata: i.Metadata,
	}
}
