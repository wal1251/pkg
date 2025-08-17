package mongo

import "go.mongodb.org/mongo-driver/bson"

type (
	Projection struct {
		value bson.M
	}
)

// NewProjection создает новый объект проекции.
func NewProjection() *Projection {
	return &Projection{value: bson.M{}}
}

// Include включает указанные поля в результат.
func (p *Projection) Include(fields ...string) *Projection {
	for _, field := range fields {
		p.value[field] = 1
	}

	return p
}

// Exclude исключает указанные поля из результата.
func (p *Projection) Exclude(fields ...string) *Projection {
	for _, field := range fields {
		p.value[field] = 0
	}

	return p
}

// Value возвращает bson.M для использования в запросах.
func (p *Projection) Value() bson.M {
	return p.value
}
