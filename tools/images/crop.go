package images

import (
	"image"

	"github.com/disintegration/imaging"
)

var (
	_ Transformer = (*CropAnchor)(nil)
	_ Transformer = (*CropCenter)(nil)
	_ Transformer = (*CropRect)(nil)
	_ Transformer = (*CropRectPercent)(nil)
)

type (
	CropAnchor struct {
		Anchor string
		Height int
		Width  int
	}

	CropCenter struct {
		Height int
		Width  int
	}

	CropRect struct {
		Left   int
		Right  int
		Top    int
		Bottom int
	}

	CropRectPercent struct {
		Left   float32
		Right  float32
		Top    float32
		Bottom float32
	}
)

func (t CropAnchor) Perform(img Image) (Image, error) {
	size := img.Bounds().Size()

	height := t.Height
	if height == 0 {
		height = size.Y
	}

	width := t.Width
	if width == 0 {
		width = size.X
	}

	if width == size.X && height == size.Y {
		return img, nil
	}

	if t.Anchor == "" {
		return img.Transform(func(img image.Image) image.Image { return imaging.CropCenter(img, width, height) }), nil
	}

	anchor, err := Anchor(t.Anchor)
	if err != nil {
		return Image{}, err
	}

	return img.Transform(func(img image.Image) image.Image { return imaging.CropAnchor(img, width, height, anchor) }), nil
}

func (t CropCenter) Perform(img Image) (Image, error) {
	size := img.Bounds().Size()

	height := t.Height
	if height == 0 {
		height = size.Y
	}

	width := t.Width
	if width == 0 {
		width = size.X
	}

	if width == size.X && height == size.Y {
		return img, nil
	}

	return img.Transform(func(img image.Image) image.Image { return imaging.CropCenter(img, width, height) }), nil
}

func (t CropRect) Perform(img Image) (Image, error) {
	if t.Left == 0 && t.Top == 0 &&
		t.Right == 0 && t.Bottom == 0 {
		return img, nil
	}

	size := img.Bounds().Size()

	return img.Transform(func(img image.Image) image.Image {
		return imaging.Crop(img, image.Rect(
			t.Left,
			t.Top,
			size.X-t.Right,
			size.Y-t.Bottom,
		))
	}), nil
}

func (t CropRectPercent) Perform(img Image) (Image, error) {
	size := img.Bounds().Size()

	return CropRect{
		percentOfEdge(size.X, t.Left),
		percentOfEdge(size.Y, t.Right),
		percentOfEdge(size.X, t.Top),
		percentOfEdge(size.Y, t.Bottom),
	}.Perform(img)
}
