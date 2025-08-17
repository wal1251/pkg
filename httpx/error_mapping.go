package httpx

import (
	"net/http"

	"github.com/wal1251/pkg/core/errs"
)

const (
	// StatusErrorDefault код статуса http, возвращаемый по-умолчанию (если тип ошибки не определен).
	StatusErrorDefault = http.StatusBadRequest

	// StatusErrorSystemFailure код статуса ошибки по-умолчанию для обще-системных сбоев.
	StatusErrorSystemFailure = http.StatusInternalServerError
)

type ErrorToStatusMapper struct {
	Mapping       map[errs.Type]int
	Default       int
	SystemFailure int
}

func (m *ErrorToStatusMapper) Status(typ errs.Type) int {
	if status, ok := m.Mapping[typ]; ok {
		return status
	}

	return StatusErrorDefault
}

func NewErrorToStatusMapper(mapping map[errs.Type]int) *ErrorToStatusMapper {
	return &ErrorToStatusMapper{
		Mapping:       mapping,
		Default:       StatusErrorDefault,
		SystemFailure: StatusErrorSystemFailure,
	}
}

func DefaultErrorToStatusMapping() map[errs.Type]int {
	return map[errs.Type]int{
		errs.TypeIllegalArgument:    http.StatusBadRequest,
		errs.TypeAuthFailure:        http.StatusUnauthorized,
		errs.TypeForbidden:          http.StatusForbidden,
		errs.TypeNotFound:           http.StatusNotFound,
		errs.TypeConflict:           http.StatusConflict,
		errs.TypeCancelled:          http.StatusGone,
		errs.TypeHasReferences:      http.StatusLocked,
		errs.TypeSystemFailure:      StatusErrorSystemFailure,
		errs.TypeNotImplemented:     http.StatusNotImplemented,
		errs.TypeServiceUnavailable: http.StatusServiceUnavailable,
		errs.TypeTooManyRequests:    http.StatusTooManyRequests,
		errs.TypeUnauthenticated:    http.StatusUnauthorized,
		errs.TypePermissionDenied:   http.StatusForbidden,
		errs.TypeResourceExhausted:  http.StatusTooManyRequests,
		errs.TypeFailedPrecondition: http.StatusPreconditionFailed,
		errs.TypeAborted:            http.StatusConflict,
		errs.TypeOutOfRange:         http.StatusRequestedRangeNotSatisfiable,
		errs.TypeUnimplemented:      http.StatusNotImplemented,
		errs.TypeInternal:           http.StatusInternalServerError,
		errs.TypeUnavailable:        http.StatusServiceUnavailable,
		errs.TypeDataLoss:           http.StatusInternalServerError,
		errs.TypeCanceled:           http.StatusGone,
	}
}
