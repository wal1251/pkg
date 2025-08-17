package images

import (
	"bytes"
	"io"

	"github.com/wal1251/pkg/tools/size"
)

var _ PFactory = (*PreviewFactory)(nil)

type (
	PreviewFactory struct {
		Options    PreviewOptions
		Name       string
		Tags       []string
		Format     *string
		Transforms []Transformer
		Inherited  []PreviewFactory
	}
)

func (f *PreviewFactory) Names() []string {
	names := make([]string, 1, len(f.Inherited)+1)
	names[0] = f.Name

	for _, inherited := range f.Inherited {
		names = append(names, inherited.Names()...)
	}

	return names
}

func (f *PreviewFactory) CreatePreview(job PreviewGeneratorJob, img Image) error {
	format := job.Format
	if f.Format != nil {
		format = *f.Format
	}

	if format == "" {
		format = "png"
	}

	preview := Preview{Name: f.Name, Format: format}
	previewIsNeeded, err := job.BeforeGenerate(&job, preview)
	if err != nil {
		return err
	}

	if !previewIsNeeded && len(f.Inherited) == 0 {
		return nil
	}

	img, err = TransformerChain(f.Transforms...).Perform(img)
	if err != nil {
		return err
	}

	if previewIsNeeded {
		buf := bytes.NewBuffer(make([]byte, 0))
		cnt := &size.CountingWriter{}
		w := io.MultiWriter(buf, cnt)

		if _, err = (Encode{
			Format:      format,
			Destination: w,
			JpegQuality: f.Options.JpegQuality,
		}.Perform(img)); err != nil {
			return err
		}

		preview.Content = bytes.NewReader(buf.Bytes())
		preview.Size = cnt.Written

		if err = job.AfterGenerate(preview); err != nil {
			return err
		}
	}

	for _, inherited := range f.Inherited {
		if err = inherited.CreatePreview(job, img); err != nil {
			return err
		}
	}

	return nil
}

func (f *PreviewFactory) AddTransform(t Transformer) {
	if f.Transforms == nil {
		f.Transforms = make([]Transformer, 0)
	}

	f.Transforms = append(f.Transforms, t)
}
