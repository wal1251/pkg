// Package annotation содержит примитивы и хелперы для аннотирования схемы данных ent.
package annotation

type CustomAnnotation struct {
	CDC         string
	Description string
}

// Name реализует интерфейс ent.Annotation.
func (CustomAnnotation) Name() string {
	return "CustomAnnotation"
}
