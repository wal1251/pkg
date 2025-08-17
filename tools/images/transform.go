package images

import (
	"errors"
	"fmt"
	"math"

	"github.com/disintegration/imaging"

	"github.com/wal1251/pkg/tools/reflection"
)

const (
	AnchorCenter      = "Center"
	AnchorTopLeft     = "TopLeft"
	AnchorTop         = "Top"
	AnchorTopRight    = "TopRight"
	AnchorLeft        = "Left"
	AnchorRight       = "Right"
	AnchorBottomLeft  = "BottomLeft"
	AnchorBottom      = "Bottom"
	AnchorBottomRight = "BottomRight"
)

const percents100 = 100

var (
	ErrUnknownResampleFilter = errors.New("unknown resample filter")
	ErrUnknownAnchor         = errors.New("unknown anchor")
)

var _ Transformer = (Transform)(nil)

type (
	Transform func(Image) (Image, error)
)

func (t Transform) Perform(i Image) (Image, error) {
	if t == nil {
		return i, nil
	}

	return t(i)
}

func TransformerChain(transformers ...Transformer) Transformer {
	return Transform(func(img Image) (Image, error) {
		var err error

		for _, transformer := range transformers {
			img, err = transformer.Perform(img)
			if err != nil {
				return Image{}, err
			}
		}

		return img, nil
	})
}

func ResampleFilter(name string) (imaging.ResampleFilter, error) {
	filter, ok := reflection.MapByTypeName(
		imaging.NearestNeighbor,
		imaging.Box,
		imaging.Linear,
		imaging.Hermite,
		imaging.MitchellNetravali,
		imaging.CatmullRom,
		imaging.BSpline,
		imaging.Gaussian,
		imaging.Bartlett,
		imaging.Lanczos,
		imaging.Hann,
		imaging.Hamming,
		imaging.Blackman,
		imaging.Cosine,
	)[name]
	if !ok {
		return imaging.ResampleFilter{}, fmt.Errorf("%w: %s", ErrUnknownResampleFilter, name)
	}

	return filter, nil
}

func Anchor(name string) (imaging.Anchor, error) {
	anchor, ok := map[string]imaging.Anchor{
		AnchorCenter:      imaging.Center,
		AnchorTopLeft:     imaging.TopLeft,
		AnchorTop:         imaging.Top,
		AnchorTopRight:    imaging.TopRight,
		AnchorLeft:        imaging.Left,
		AnchorRight:       imaging.Right,
		AnchorBottomLeft:  imaging.BottomLeft,
		AnchorBottom:      imaging.Bottom,
		AnchorBottomRight: imaging.BottomRight,
	}[name]
	if !ok {
		return imaging.Anchor(0), fmt.Errorf("%w: %s", ErrUnknownAnchor, name)
	}

	return anchor, nil
}

func DefaultTransformerFactory() TransformerFactory {
	factory := make(TransformerFactory)
	factory.Register((*Fit)(nil))
	factory.Register((*CropRect)(nil))
	factory.Register((*CropRectPercent)(nil))
	factory.Register((*CropAnchor)(nil))
	factory.Register((*CropCenter)(nil))
	factory.Register((*Encode)(nil))

	return factory
}

func percentOfEdge(edge int, percent float32) int {
	return int(math.Round(float64(float32(edge) * percent / percents100)))
}
