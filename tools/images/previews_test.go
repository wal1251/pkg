package images_test

import (
	"image"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/tools/images"
)

func TestPreviewGenerator_Generate_preview_selector(t *testing.T) {
	gen := NewPreviewsGenerator(t, `{
	"preview1": {
		"tags": ["preview_a"]
	},
	"preview2": {
		"tags": ["preview_a"]
	},
	"preview3": {
		"tags": ["preview_b"],
		"inherited": {
			"preview31": {
				"inherited": {
					"preview311": {},
					"preview312": {}
				}
			},
			"preview32": {}
		}
	},
	"preview4": {
		"tags": ["preview_b"]
	}
}`, nil)

	tests := []struct {
		name         string
		job          images.PreviewGeneratorJob
		filter       func(preview images.Preview) bool
		wantPreviews []string
	}{
		{
			name:         "Базовый кейс",
			job:          images.PreviewGeneratorJob{Tag: "preview_a"},
			wantPreviews: []string{"preview1", "preview2"},
		},
		{
			name:         "Выборка preview в том числе с наследуемыми",
			job:          images.PreviewGeneratorJob{Tag: "preview_b"},
			wantPreviews: []string{"preview3", "preview31", "preview311", "preview312", "preview32", "preview4"},
		},
		{
			name: "Выборка preview в том числе с наследуемыми",
			job:  images.PreviewGeneratorJob{Tag: "preview_b"},
			filter: func(preview images.Preview) bool {
				return preview.Name == "preview312" || preview.Name == "preview4"
			},
			wantPreviews: []string{"preview312", "preview4"},
		},
	}

	for _, tt := range tests {
		sequence := make([]string, 0)

		t.Run(tt.name, func(t *testing.T) {
			job := tt.job.
				OnBeforeGenerate(func(_ *images.PreviewGeneratorJob, preview images.Preview) (bool, error) {
					return tt.filter == nil || tt.filter(preview), nil
				}).
				OnAfterGenerate(func(preview images.Preview) error {
					sequence = append(sequence, preview.Name)
					return nil
				})
			require.NoError(t, gen.Generate(job, NewBlankImage(0, 0, 100, 100)))
			assert.ElementsMatch(t, tt.wantPreviews, sequence)
		})
	}
}

func TestPreviewGenerator_Generate_format_selector(t *testing.T) {
	gen := NewPreviewsGenerator(t, `{
	"jpegPreview": {
		"tags": ["jpegWanted"],
		"format": "jpeg"
	},
	"defaultPreview": {
		"tags": ["defaultWanted"]
	}
}`, nil)

	tests := []struct {
		name        string
		job         images.PreviewGeneratorJob
		wantFormat  string
		wantContent string
	}{
		{
			name:        "Если формат неизвестен, сохранит в png",
			job:         images.PreviewGeneratorJob{Tag: "defaultWanted"},
			wantFormat:  "png",
			wantContent: "image/png",
		},
		{
			name:        "Если целевой формат не задан, сохранит в исходном",
			job:         images.PreviewGeneratorJob{Tag: "defaultWanted", Format: "jpeg"},
			wantFormat:  "jpeg",
			wantContent: "image/jpeg",
		},
		{
			name:        "Если целевой формат задан, сохранит в целевом",
			job:         images.PreviewGeneratorJob{Tag: "jpegWanted"},
			wantFormat:  "jpeg",
			wantContent: "image/jpeg",
		},
		{
			name:        "Если целевой формат задан, сохранит в целевом, даже если исходный задан",
			job:         images.PreviewGeneratorJob{Tag: "jpegWanted", Format: "bmp"},
			wantFormat:  "jpeg",
			wantContent: "image/jpeg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := tt.job.
				OnBeforeGenerate(func(_ *images.PreviewGeneratorJob, preview images.Preview) (bool, error) {
					assert.Equal(t, tt.wantFormat, preview.Format)
					return true, nil
				}).
				OnAfterGenerate(func(preview images.Preview) error {
					assert.Equal(t, tt.wantFormat, preview.Format)
					assert.Equal(t, tt.wantContent, SniffFormat(t, preview.Content))
					return nil
				})
			require.NoError(t, gen.Generate(job, NewBlankImage(0, 0, 100, 100)))
		})
	}
}

func TestPreviewGenerator_Generate_basic_transformation(t *testing.T) {
	gen := NewPreviewsGenerator(t, `
{
	"main": {
		"tags": [ "foo" ],
		"format": "jpeg",
		"transforms": [
			{
				"type": "images.CropCenter",
				"properties": {
					"Height": 100,
					"Width": 175
				}
			}
		],
		"inherited": {
			"secondary1": {
				"format": "png",
				"transforms": [
					{
						"type": "images.CropCenter",
						"properties": {
							"Height": 50,
							"Width": 100
						}
					}
				]
			}
		}
	}
}`, nil)

	job := images.PreviewGeneratorJob{Tag: "foo", Format: "bmp"}.
		OnBeforeGenerate(func(_ *images.PreviewGeneratorJob, preview images.Preview) (bool, error) {
			switch preview.Name {
			case "main":
				assert.Equal(t, "jpeg", preview.Format)
			case "secondary":
				assert.Equal(t, "png", preview.Format)
			}
			return true, nil
		}).
		OnAfterGenerate(func(preview images.Preview) error {
			img := images.MakeImage()
			err := img.Decode(preview.Content)
			switch preview.Name {
			case "main":
				assert.Equal(t, "jpeg", preview.Format)
				assert.Equal(t, 175, img.Bounds().Size().X)
				assert.Equal(t, 100, img.Bounds().Size().Y)
			case "secondary":
				assert.Equal(t, "png", preview.Format)
				assert.Equal(t, 100, img.Bounds().Size().X)
				assert.Equal(t, 50, img.Bounds().Size().Y)
			}
			return err
		})
	assert.NoError(t, gen.Generate(job, NewBlankImage(0, 0, 200, 200)))
}

func TestPreviewGenerator_Generate_conditional_transformation(t *testing.T) {
	gen := NewPreviewsGenerator(t, `{
	"basic": {
		"tags": [ "foo" ],
		"transforms": [
			{
				"type": "images.CropCenter",
				"conditions": [ "OrientationHorizontal" ],
				"properties": {
					"Height": 50,
					"Width": 175
				}
			},
			{
				"type": "images.CropCenter",
				"conditions": [ "OrientationVertical" ],
				"properties": {
					"Height": 100,
					"Width": 100
				}
			}
		]
	}
}`, nil)

	tests := []struct {
		name       string
		job        images.PreviewGeneratorJob
		height     int
		width      int
		wantHeight int
		wantWidth  int
	}{
		{
			name:       "Должна сработать трансформация для горизонтального изображения",
			job:        images.PreviewGeneratorJob{Tag: "foo"},
			width:      200,
			height:     190,
			wantWidth:  175,
			wantHeight: 50,
		},
		{
			name:       "Должна сработать трансформация для вертикального изображения",
			job:        images.PreviewGeneratorJob{Tag: "foo"},
			width:      190,
			height:     200,
			wantWidth:  100,
			wantHeight: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := tt.job.
				OnBeforeGenerate(func(_ *images.PreviewGeneratorJob, preview images.Preview) (bool, error) { return true, nil }).
				OnAfterGenerate(func(preview images.Preview) error {
					img := images.MakeImage()
					require.NoError(t, img.Decode(preview.Content))

					assert.Equal(t, tt.wantHeight, img.Bounds().Size().Y)
					assert.Equal(t, tt.wantWidth, img.Bounds().Size().X)
					return nil
				})
			require.NoError(t, gen.Generate(job, NewBlankImage(0, 0, tt.width, tt.height)))
		})
	}
}

func TestPreviewGenerator_Generate_tag_not_found(t *testing.T) {
	gen := NewPreviewsGenerator(t, `
{
	"main": {
		"tags": [ "foo" ]
	}
}`, nil)

	require.Error(t, gen.Generate(images.PreviewGeneratorJob{Tag: "baz"},
		NewBlankImage(0, 0, 100, 200)),
	)
}

func NewBlankImage(x0, y0, x1, y1 int) images.Image {
	img := images.MakeImage()
	img.Image = image.NewNRGBA(image.Rect(x0, y0, x1, y1))
	return img
}

func NewPreviewsGenerator(t *testing.T, json string, tune func(images.TransformerFactory, images.ConditionsRepo)) *images.PreviewGenerator {
	factory := images.DefaultTransformerFactory()
	conditions := images.DefaultConditions()
	if tune != nil {
		tune(factory, conditions)
	}
	previewGenerator, err := images.NewPreviewGenerator(func() ([]images.PreviewFactory, error) {
		return images.LoadPreviewsSchemaFromJSON(
			[]byte(json),
			factory,
			images.DefaultConditions(),
			images.PreviewOptions{})
	})
	require.NoError(t, err)
	return previewGenerator
}

func SniffFormat(t *testing.T, img io.Reader) string {
	data, err := io.ReadAll(img)
	require.NoError(t, err)
	return http.DetectContentType(data)
}
