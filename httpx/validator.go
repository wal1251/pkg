package httpx

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"

	"github.com/wal1251/pkg/core/logs"
)

var (
	_ RequestValidator            = (*OpenAPIRequestValidator)(nil)
	_ error                       = (*ValidationError)(nil)
	_ interface{ Unwrap() error } = (*ValidationError)(nil)
)

// TODO: Создать универсальный интерфейс валидации, open api валидатор, как отдельная реализация.

type (
	// RequestValidator предназначен для валидации входящих запросов к серверу.
	RequestValidator interface {
		// Validate выполняет валидацию запроса r, если запрос валиден, то возвращает nil, в противном случае вернет
		// ErrValidationFailed, если в работе RequestValidator происходили иные ошибки, вернет error отличный
		// от ErrValidationFailed.
		Validate(r *http.Request) error
	}

	// RequestBodyValidationPredicate предикат валидации тела запроса, если вернет true валидация тела будет выполняться,
	// false - не будет. Предикат вызывается валидатором перед осуществлением валидации каждого запроса.
	RequestBodyValidationPredicate func(*http.Request, *routers.Route) bool

	// RequestParameter тип проверяемого параметра запроса.
	RequestParameter string

	// ValidationError объект, содержащий информацию об ошибке валидации.
	ValidationError struct {
		Err       error
		Message   string
		Field     string
		Parameter RequestParameter
		Value     any
	}

	// OpenAPIRequestValidator валидирует запрос к серверу согласно схеме OpenAPI.
	OpenAPIRequestValidator struct {
		router                  routers.Router
		bodyValidationPredicate RequestBodyValidationPredicate
	}
)

// Validate выполняет валидацию запроса r согласно спецификации OpenAPI заданной для роутера (см. routers.Router) в
// OpenAPIRequestValidator. Если запрос валиден, возвращает nil, в противном случае ErrValidationFailed, если в работе
// валидатора происходили иные ошибки, вернет error отличный от ErrValidationFailed.
func (v *OpenAPIRequestValidator) Validate(request *http.Request) error {
	logger := logs.FromContext(request.Context())

	route, pathParams, err := v.router.FindRoute(request)
	if err != nil {
		return fmt.Errorf("failed to get route from api schema: %w", err)
	}

	options := &openapi3filter.Options{
		AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		ExcludeRequestBody: !v.bodyValidationPredicate.Accept(request, route),
	}

	logger.Debug().Msgf("body validation against api schema disabled: %v", options.ExcludeRequestBody)

	requestValidationInput := &openapi3filter.RequestValidationInput{
		Route:      route,
		PathParams: pathParams,
		Request:    request,
		Options:    options,
	}

	return validationErrorFromRequestError(openapi3filter.ValidateRequest(request.Context(), requestValidationInput))
}

func (p RequestBodyValidationPredicate) Accept(r *http.Request, route *routers.Route) bool {
	if p != nil {
		return p(r, route)
	}

	return true
}

func (p RequestBodyValidationPredicate) And(predicates ...RequestBodyValidationPredicate) RequestBodyValidationPredicate {
	final := p
	for _, predicate := range predicates {
		current := final
		next := predicate
		final = func(r *http.Request, route *routers.Route) bool {
			if current(r, route) {
				return next(r, route)
			}

			return false
		}
	}

	return final
}

func (p RequestParameter) IsQuery() bool {
	return p == "query"
}

func (p RequestParameter) IsPath() bool {
	return p == "path"
}

func (v *ValidationError) Error() string {
	return v.Message
}

func (v *ValidationError) Unwrap() error {
	return v.Err
}

func (v *ValidationError) Is(err error) bool {
	var validationErr *ValidationError

	return errors.As(err, &validationErr)
}

// setSchemaError Устанавливает ошибку валидации в структуру ValidationError.
func (v *ValidationError) setSchemaError(schemaError *openapi3.SchemaError) { //nolint:unused
	var requestError *openapi3filter.RequestError

	reason := schemaError.Reason
	if schemaError.SchemaField == "pattern" {
		reason = "string doesn't match the regular expression"
	}

	v.Err = schemaError
	v.Message = reason
	v.Value = schemaError.Value

	if !errors.As(schemaError, &requestError) {
		return
	}

	if requestError.Parameter == nil {
		v.Field = strings.Join(schemaError.JSONPointer(), ".")

		return
	}

	v.Field = requestError.Parameter.Name
	v.Parameter = RequestParameter(requestError.Parameter.In)
}

// setSchemaErrorWithField Устанавливает ошибку валидации в структуру c полями которые не хватают ValidationError.
func (v *ValidationError) setSchemaErrorWithField(schemaError *openapi3.SchemaError) {
	var requestError *openapi3filter.RequestError

	reason := schemaError.Reason
	if schemaError.SchemaField == "pattern" {
		reason = "string doesn't match the regular expression"
	}

	pathElements := schemaError.JSONPointer()
	field := ""
	if len(pathElements) > 0 {
		field = pathElements[len(pathElements)-1]
	}

	v.Err = schemaError
	v.Message = reason
	v.Value = schemaError.Value
	v.Field = field

	if !errors.As(schemaError, &requestError) {
		return
	}

	v.Err = schemaError
	v.Message = reason
	v.Value = schemaError.Value

	if !errors.As(schemaError, &requestError) {
		return
	}

	if requestError.Parameter == nil {
		v.Field = strings.Join(schemaError.JSONPointer(), ".")

		return
	}

	v.Field = requestError.Parameter.Name
	v.Parameter = RequestParameter(requestError.Parameter.In)
}

func ValidateRequestBody() RequestBodyValidationPredicate {
	return func(_ *http.Request, route *routers.Route) bool {
		return route.Operation.RequestBody != nil && route.Operation.RequestBody.Value != nil
	}
}

func SchemaTraverse(schemaRef *openapi3.SchemaRef, consumer func(schema *openapi3.Schema) bool) bool {
	if schemaRef == nil || schemaRef.Value == nil {
		return true
	}

	if !consumer(schemaRef.Value) {
		return false
	}

	if !SchemaTraverse(schemaRef.Value.Items, consumer) {
		return false
	}

	for _, property := range schemaRef.Value.Properties {
		if !SchemaTraverse(property, consumer) {
			return false
		}
	}

	return true
}

func DontValidateBinaryBody() RequestBodyValidationPredicate {
	return func(_ *http.Request, route *routers.Route) bool {
		for _, content := range route.Operation.RequestBody.Value.Content {
			if !SchemaTraverse(
				content.Schema,
				func(schema *openapi3.Schema) bool { return schema.Format != "binary" },
			) {
				return false
			}
		}

		return true
	}
}

func DontValidateXMLBody() RequestBodyValidationPredicate {
	return func(r *http.Request, _ *routers.Route) bool {
		return r.Header.Get("Content-Type") != ContentTypeXML
	}
}

func NewFieldValidationError(err error, field string) error {
	return &ValidationError{
		Err:     err,
		Message: err.Error(),
		Field:   field,
	}
}

func validationErrorFromRequestError(err error) error {
	var schemaError *openapi3.SchemaError
	var requestError *openapi3filter.RequestError

	if err == nil {
		return nil
	}

	if errors.As(err, &schemaError) {
		validationError := &ValidationError{}
		validationError.setSchemaErrorWithField(schemaError)

		return validationError
	}

	if errors.As(err, &requestError) {
		msg := requestError.Reason
		if msg == "" {
			msg = requestError.Error()
		}

		validationError := &ValidationError{
			Err:     requestError.Err,
			Message: msg,
		}

		if requestError.Parameter != nil {
			validationError.Field = requestError.Parameter.Name
			validationError.Parameter = RequestParameter(requestError.Parameter.In)
		}

		return validationError
	}

	return err
}

func ValidatorError(err error) (*ValidationError, bool) {
	if err != nil {
		var validationError *ValidationError

		if errors.As(err, &validationError) {
			return validationError, true
		}
	}

	return nil, false
}

func IsParseError(err error) bool {
	if err != nil {
		var validationError *openapi3filter.ParseError

		return errors.As(err, &validationError)
	}

	return false
}

func NewRequestValidator(
	openAPIRouter routers.Router,
	predicates ...RequestBodyValidationPredicate,
) *OpenAPIRequestValidator {
	return &OpenAPIRequestValidator{
		router:                  openAPIRouter,
		bodyValidationPredicate: ValidateRequestBody().And(predicates...),
	}
}
