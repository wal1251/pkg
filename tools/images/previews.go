package images

import (
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/wal1251/pkg/tools/anyobj"
	"github.com/wal1251/pkg/tools/collections"
	"github.com/wal1251/pkg/tools/serial"
)

var (
	ErrPreviewSchemaNotFound    = errors.New("preview schema not found")
	ErrPreviewSchemaParseFailed = errors.New("preview schema parse failed")
)

func LoadPreviewsSchemaFromFileYAML(fileName string, transformers TransformerFactory, conditions ConditionsRepo,
	options PreviewOptions,
) ([]PreviewFactory, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to load previews creation schema: %w", err)
	}

	yamlFile, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to load previews creation schema: %w", err)
	}

	return LoadPreviewsSchemaFromYAML(yamlFile, transformers, conditions, options)
}

func LoadPreviewsSchemaFromYAML(yamlFile []byte, transformers TransformerFactory, conditions ConditionsRepo, options PreviewOptions) ([]PreviewFactory, error) {
	schema := make(map[string]any)
	if err := yaml.Unmarshal(yamlFile, &schema); err != nil {
		return nil, fmt.Errorf("can't unmarshal preview schema from yaml: %w", err)
	}

	return ParsePreviewsSchema(schema, transformers, conditions, options)
}

func LoadPreviewsSchemaFromJSON(jsonFile []byte, transformers TransformerFactory, conditions ConditionsRepo, options PreviewOptions) ([]PreviewFactory, error) {
	schema, err := serial.FromBytes(jsonFile, serial.JSONDecode[map[string]any])
	if err != nil {
		return nil, err
	}

	return ParsePreviewsSchema(schema, transformers, conditions, options)
}

func ParsePreviewsSchema(schema map[string]any, transformers TransformerFactory, conditions ConditionsRepo, //nolint: gocognit
	options PreviewOptions,
) ([]PreviewFactory, error) {
	previewFactories := make([]PreviewFactory, 0)

	for name, preview := range schema {
		previewSchema, ok := preview.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("%w: preview schema object is expected (%s)", ErrPreviewSchemaParseFailed, name)
		}

		previewFactory, err := parsePreviewSchema(name, previewSchema)
		if err != nil {
			return nil, err
		}

		if transforms, ok := previewSchema["transforms"]; ok {
			transformsSlice, ok := transforms.([]any)
			if !ok {
				return nil, fmt.Errorf("%w: transform schema object is expected (%s)", ErrPreviewSchemaParseFailed, name)
			}

			for index, transform := range transformsSlice {
				transformSchema, ok := transform.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("%w: object is expected as transform schema element: %s %d",
						ErrPreviewSchemaParseFailed, name, index)
				}

				transformer, err := parsePreviewTransformSchema(transformSchema, transformers, conditions)
				if err != nil {
					return nil, fmt.Errorf("failed to parse transform of %s: %w", name, err)
				}

				previewFactory.AddTransform(transformer)
			}
		}

		if inherited, ok := previewSchema["inherited"]; ok {
			inheritedSchema, ok := inherited.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("%w: object is expected as inherited schema: %s",
					ErrPreviewSchemaParseFailed, name)
			}

			inheritedPreviews, err := ParsePreviewsSchema(inheritedSchema, transformers, conditions, options)
			if err != nil {
				return nil, fmt.Errorf("failed to parse inherited schema of %s: %w", name, err)
			}

			if len(inheritedPreviews) != 0 {
				previewFactory.Inherited = inheritedPreviews
			}
		}

		previewFactories = append(previewFactories, previewFactory)
	}

	return previewFactories, nil
}

func parsePreviewSchema(name string, schema map[string]any) (PreviewFactory, error) {
	factory := PreviewFactory{Name: name}
	if format, ok := schema["format"]; ok {
		formatString, ok := format.(string)
		if !ok {
			return PreviewFactory{}, fmt.Errorf("%w: string expected as format property: %s",
				ErrPreviewSchemaParseFailed, name)
		}

		factory.Format = &formatString
	}

	if tags, ok := schema["tags"]; ok {
		tagsAny, ok := tags.([]any)
		if !ok {
			return PreviewFactory{}, fmt.Errorf("%w: array expected as tags property (%s)",
				ErrPreviewSchemaParseFailed, name)
		}

		for index, tag := range tagsAny {
			tagString, ok := tag.(string)
			if !ok {
				return PreviewFactory{}, fmt.Errorf("%w: string expected as tags array element: %s %d",
					ErrPreviewSchemaParseFailed, name, index)
			}

			if factory.Tags == nil {
				factory.Tags = make([]string, 0, 1)
			}

			factory.Tags = append(factory.Tags, tagString)
		}
	}

	return factory, nil
}

func parsePreviewTransformSchema(schema map[string]any, transformers TransformerFactory, knownConditions ConditionsRepo) (Transformer, error) {
	name, ok := schema["type"].(string)
	if !ok {
		return nil, fmt.Errorf("%w: string is expected as transformer type", ErrPreviewSchemaParseFailed)
	}

	transformer := transformers.Create(name)
	if transformer == nil {
		return nil, fmt.Errorf("%w: unknown transformer type: %s", ErrPreviewSchemaParseFailed, name)
	}

	if properties, ok := schema["properties"]; ok {
		if err := anyobj.SafeCopy(properties, &transformer); err != nil {
			return nil, fmt.Errorf("failed to parse %s transformer properties: %w", name, err)
		}
	}

	if conditions, ok := schema["conditions"]; ok {
		conditionsSlice, ok := conditions.([]any)
		if !ok {
			return nil, fmt.Errorf("%w: array is expected as conditions of transform: %s",
				ErrPreviewSchemaParseFailed, name)
		}

		if len(conditionsSlice) != 0 {
			predicates := collections.Map(conditionsSlice, func(value any) Predicate {
				if condition, ok := value.(string); ok {
					return knownConditions.Predicate(condition)
				}

				return (Predicate)(nil)
			})

			transformer = MakeCondition(predicates...).MakeTransformer(transformer)
		}
	}

	return transformer, nil
}
