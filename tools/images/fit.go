package images

import (
	"image"

	"github.com/disintegration/imaging"
)

var _ Transformer = (*Fit)(nil)

type (
	Fit struct {
		ResampleFilter string
		Height         int
		Width          int
	}
)

func (t Fit) Perform(img Image) (Image, error) {
	size := img.Bounds().Size()

	width := t.Width
	if width == 0 {
		width = size.X
	}

	height := t.Height
	if height == 0 {
		height = size.Y
	}

	if width == size.X && height == size.Y {
		return img, nil
	}

	if t.ResampleFilter == "" {
		return img.Transform(func(img image.Image) image.Image { return imaging.Fit(img, width, height, imaging.Lanczos) }), nil
	}

	filter, err := ResampleFilter(t.ResampleFilter)
	if err != nil {
		return Image{}, err
	}

	return img.Transform(func(img image.Image) image.Image { return imaging.Fit(img, width, height, filter) }), nil
}
