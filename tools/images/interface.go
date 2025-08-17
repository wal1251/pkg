package images

import (
	"image"
	"io"
)

type (
	Condition interface {
		Check(Image) bool
		MakeTransformer(Transformer) Transformer
	}

	ConditionRepo interface {
		Predicate(name string) Predicate
	}

	Transformer interface {
		Perform(Image) (Image, error)
	}

	TransformerRepo interface {
		Register(v any)
		Create(name string) Transformer
	}

	Imageable interface {
		Decode(input io.Reader) error
		Transform(t func(image.Image) image.Image) Image
	}

	PGeneratorJob interface {
		BeforeGenerate(job *PreviewGeneratorJob, preview Preview) (bool, error)
		AfterGenerate(preview Preview) error
		OnBeforeGenerate(h func(*PreviewGeneratorJob, Preview) (bool, error)) PreviewGeneratorJob
		OnAfterGenerate(h func(Preview) error) PreviewGeneratorJob
	}

	PGenerator interface {
		Generate(job PreviewGeneratorJob, img Image) error
		Names(job PreviewGeneratorJob) []string
	}

	PFactory interface {
		Names() []string
		CreatePreview(job PreviewGeneratorJob, img Image) error
		AddTransform(t Transformer)
	}
)
