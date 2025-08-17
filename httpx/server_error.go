package httpx

import (
	"errors"

	"github.com/wal1251/pkg/core"
	"github.com/wal1251/pkg/core/errs"
)

var _ core.Map[error, ServerError] = MakeServerError

type ServerError struct {
	Code    string            `json:"code,omitempty"`
	Message string            `json:"message,omitempty"`
	Error   string            `json:"error,omitempty"`
	Details map[string]string `json:"details,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func MakeServerError(err error) ServerError {
	return ServerError{
		Code:    errs.AsReason(err).Code,
		Message: err.Error(),
	}
}

func MakeSecureServerError(err error) ServerError {
	errNum := errs.AsReason(err).ErrNum
	if errNum == "" || errNum == "0.0" || errNum == "0" {
		return ServerError{
			Error: "1.0",
		}
	}
	details := errs.AsReason(err).Details
	if details != nil {
		return ServerError{
			Error:   errNum,
			Details: details,
		}
	}

	return ServerError{
		Error: errNum,
	}
}

func MakeSecureValidatorServerError(err error) ServerError {
	var serverError *errs.WrappingError
	// make it ValidatorError
	if errors.As(err, &serverError) {
		if len(serverError.Fields) != 0 {
			return ServerError{
				Error:  "1.0",
				Fields: serverError.Fields,
			}
		}
	}

	return ServerError{
		Error: "1.0",
	}
}
