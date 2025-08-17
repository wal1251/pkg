package search

import (
	"fmt"
	"strings"

	"github.com/wal1251/pkg/providers/es/api"
	"github.com/wal1251/pkg/tools/collections"
	"github.com/wal1251/pkg/tools/serial"
)

const (
	QueryBoolKey        = "bool"
	QueryBoolMustKey    = "must"
	QueryBoolFilterKey  = "filter"
	QueryBoolShouldKey  = "should"
	QueryBoolMustNotKey = "must_not"
)

var _ Query = (*QueryBoolean)(nil)

type (
	// QueryBoolean запрос состоящий из серии логических предложений.
	QueryBoolean struct {
		Must    Queries
		Filter  Queries
		Should  Queries
		MustNot Queries
	}
)

func (q *QueryBoolean) AddMust(queries ...Query) {
	if q.Must == nil {
		q.Must = make(Queries, 0, len(queries))
	}
	q.Must = append(q.Must, queries...)
}

func (q *QueryBoolean) AddShould(queries ...Query) {
	if q.Should == nil {
		q.Should = make(Queries, 0, len(queries))
	}
	q.Should = append(q.Should, queries...)
}

func (q *QueryBoolean) AddFilter(queries ...Query) {
	if q.Filter == nil {
		q.Filter = make(Queries, 0, len(queries))
	}
	q.Filter = append(q.Filter, queries...)
}

func (q *QueryBoolean) AddMustNot(queries ...Query) {
	if q.MustNot == nil {
		q.MustNot = make(Queries, 0, len(queries))
	}
	q.MustNot = append(q.MustNot, queries...)
}

func (q *QueryBoolean) IsEmpty() bool {
	return len(q.Must) == 0 &&
		len(q.Should) == 0 &&
		len(q.Filter) == 0 &&
		len(q.MustNot) == 0
}

func (q *QueryBoolean) MarshalJSON() ([]byte, error) {
	return serial.ToBytes(api.Object{
		q.QueryType(): api.CreateObject(
			api.PropertyOmitempty(QueryBoolMustKey, q.Must),
			api.PropertyOmitempty(QueryBoolFilterKey, q.Filter),
			api.PropertyOmitempty(QueryBoolShouldKey, q.Should),
			api.PropertyOmitempty(QueryBoolMustNotKey, q.MustNot),
		),
	}, serial.JSONEncode[api.Object])
}

func (q *QueryBoolean) QueryType() string {
	return QueryBoolKey
}

func (q *QueryBoolean) Visitor(fields []DocumentFieldDescriptor, language string) func(any, QueryFieldDescriptor) {
	queryableFields := queryBooleanQueryableFields(fields, language)
	filterableFields := queryBooleanFilterableFields(fields)

	return func(query any, meta QueryFieldDescriptor) {
		queryString := fmt.Sprint(query)

		switch meta.Type {
		case QueryFieldTypeQuery:
			q.appendQuery(strings.TrimSpace(queryString), collections.Filter(queryableFields,
				func(d DocumentFieldDescriptor) bool { return meta.Name == "" || meta.Name == d.Name }))
		case QueryFieldTypeFilter:
			q.appendFilter(queryString, collections.Filter(filterableFields,
				func(d DocumentFieldDescriptor) bool {
					return meta.Name == d.Field ||
						(meta.Name != "" && meta.Name == d.Name) ||
						(meta.Name == "" && meta.Field == d.Field)
				}),
			)
		}
	}
}

func (q *QueryBoolean) appendQuery(query string, fields []DocumentFieldDescriptor) {
	if query == "" {
		return
	}

	queryBooleanAppend(query, q.AddMust, QueryBoolMustKey, fields)
	queryBooleanAppend(query, q.AddShould, QueryBoolShouldKey, fields)
}

func (q *QueryBoolean) appendFilter(term string, fields []DocumentFieldDescriptor) {
	for _, field := range fields {
		q.AddFilter(&QueryTerm{Field: field.Field, Value: term})
	}
}

func NewQueryBoolean(search any, language string, descriptor QueryDescriptor, document []DocumentFieldDescriptor, opts ...func(*QueryBoolean)) *QueryBoolean {
	var query QueryBoolean

	descriptor.Visit(search, query.Visitor(document, language))

	for _, opt := range opts {
		opt(&query)
	}

	return &query
}

func queryBooleanFilterableFields(fields []DocumentFieldDescriptor) []DocumentFieldDescriptor {
	return collections.Filter(fields, func(d DocumentFieldDescriptor) bool { return d.Type == QueryBoolFilterKey })
}

func queryBooleanQueryableFields(fields []DocumentFieldDescriptor, language string) []DocumentFieldDescriptor {
	return collections.Filter(fields, func(d DocumentFieldDescriptor) bool {
		return (strings.HasPrefix(d.Type, QueryBoolMustKey) || strings.HasPrefix(d.Type, QueryBoolShouldKey)) &&
			(d.Language == nil || *d.Language == language)
	})
}

func queryBooleanSubQueries(fields []DocumentFieldDescriptor, subQuery, defaultQuery string) map[string][]DocumentFieldDescriptor {
	fields = collections.Filter(fields, func(f DocumentFieldDescriptor) bool {
		return strings.HasPrefix(f.Type, subQuery)
	})

	return collections.Group(fields, func(d DocumentFieldDescriptor) string {
		if strings.Contains(d.Type, ".") {
			return strings.Split(d.Type, ".")[1]
		}

		return defaultQuery
	})
}

func queryBooleanAppend(query string, booleanAdd func(...Query), booleanType string, fields []DocumentFieldDescriptor) {
	bySearchTypes := queryBooleanSubQueries(fields, booleanType, QueryMultiMatchKey)
	if len(bySearchTypes) == 0 {
		return
	}

	if len(bySearchTypes) > 1 {
		booleanAdd(NewQueryDisjunctionMax(
			collections.Join(
				collections.Map(bySearchTypes[QueryTermKey], func(d DocumentFieldDescriptor) Query {
					return NewQueryTerm(d.Field, query)
				}),
				collections.Single[Query](NewQueryMultiMatch(query, bySearchTypes[QueryMultiMatchKey]...)),
			),
		))
	}

	for queryType, fieldsGroup := range bySearchTypes {
		switch queryType {
		case QueryTermKey:
			if len(fieldsGroup) == 1 {
				booleanAdd(NewQueryTerm(fieldsGroup[0].Field, query))
			} else {
				booleanAdd(NewQueryDisjunctionMax(
					collections.Map(fieldsGroup, func(d DocumentFieldDescriptor) Query {
						return NewQueryTerm(d.Field, query)
					}),
				))
			}
		case QueryMultiMatchKey:
			booleanAdd(NewQueryMultiMatch(query, fieldsGroup...))
		}
	}
}
