package images

import "fmt"

var _ PGenerator = (*PreviewGenerator)(nil)

type (
	PreviewGenerator struct {
		previewsByTags map[string][]*PreviewFactory
	}
)

func NewPreviewGenerator(loader func() ([]PreviewFactory, error)) (*PreviewGenerator, error) {
	factories, err := loader()
	if err != nil {
		return nil, err
	}

	tagged := make(map[string][]*PreviewFactory)
	for index, factory := range factories {
		for _, tag := range factory.Tags {
			if _, ok := tagged[tag]; !ok {
				tagged[tag] = make([]*PreviewFactory, 0, 1)
			}

			tagged[tag] = append(tagged[tag], &factories[index])
		}
	}

	return &PreviewGenerator{
		previewsByTags: tagged,
	}, err
}

func (g *PreviewGenerator) Generate(job PreviewGeneratorJob, img Image) error {
	factories, ok := g.previewsByTags[job.Tag]
	if !ok {
		return fmt.Errorf("%w: tag %s", ErrPreviewSchemaNotFound, job.Tag)
	}

	for _, factory := range factories {
		if err := factory.CreatePreview(job, img); err != nil {
			return err
		}
	}

	return nil
}

func (g *PreviewGenerator) Names(job PreviewGeneratorJob) []string {
	factories := g.previewsByTags[job.Tag]
	names := make([]string, 0, len(factories)+1)

	for _, factory := range factories {
		names = append(names, factory.Names()...)
	}

	return names
}
