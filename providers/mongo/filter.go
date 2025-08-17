package mongo

import "go.mongodb.org/mongo-driver/bson"

type (
	// Filter оборачивает bson.D для инкапсуляции.
	Filter struct {
		value bson.D
	}
)

// NewFilter создает новый пустой фильтр.
func NewFilter() *Filter {
	return &Filter{value: bson.D{}}
}

// Eq добавляет условие "равно" (field = value).
func (f *Filter) Eq(field string, value interface{}) *Filter {
	f.value = append(f.value, bson.E{Key: field, Value: value})

	return f
}

// Gt добавляет условие "больше" (field > value).
func (f *Filter) Gt(field string, value interface{}) *Filter {
	f.value = append(f.value, bson.E{Key: field, Value: bson.M{"$gt": value}})

	return f
}

// Lt добавляет условие "меньше" (field < value).
func (f *Filter) Lt(field string, value interface{}) *Filter {
	f.value = append(f.value, bson.E{Key: field, Value: bson.M{"$lt": value}})

	return f
}

// And добавляет условие "и" (логическое $and).
func (f *Filter) And(filters ...*Filter) *Filter {
	var andConditions bson.A
	for _, filter := range filters {
		andConditions = append(andConditions, filter.value)
	}
	f.value = append(f.value, bson.E{Key: "$and", Value: andConditions})

	return f
}

// Or добавляет условие "или" (логическое $or).
func (f *Filter) Or(filters ...*Filter) *Filter {
	var orConditions bson.A
	for _, filter := range filters {
		orConditions = append(orConditions, filter.value)
	}
	f.value = append(f.value, bson.E{Key: "$or", Value: orConditions})

	return f
}

// Value возвращает внутренний bson.D.
func (f *Filter) Value() bson.D {
	return f.value
}
