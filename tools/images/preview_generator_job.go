package images

var _ PGeneratorJob = (*PreviewGeneratorJob)(nil)

type (
	PreviewGeneratorJob struct {
		Tag            string
		Format         string
		beforeGenerate func(*PreviewGeneratorJob, Preview) (bool, error)
		afterGenerate  func(Preview) error
	}
)

func (c PreviewGeneratorJob) BeforeGenerate(job *PreviewGeneratorJob, preview Preview) (bool, error) {
	if c.beforeGenerate == nil {
		return true, nil
	}

	return c.beforeGenerate(job, preview)
}

func (c PreviewGeneratorJob) AfterGenerate(preview Preview) error {
	if c.afterGenerate == nil {
		return nil
	}

	return c.afterGenerate(preview)
}

func (c PreviewGeneratorJob) OnBeforeGenerate(h func(*PreviewGeneratorJob, Preview) (bool, error)) PreviewGeneratorJob {
	c.beforeGenerate = h

	return c
}

func (c PreviewGeneratorJob) OnAfterGenerate(h func(Preview) error) PreviewGeneratorJob {
	c.afterGenerate = h

	return c
}
