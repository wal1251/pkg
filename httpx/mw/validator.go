package mw

import (
	"errors"
	"net/http"

	"github.com/wal1251/pkg/core/errs"
	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/httpx"
)

func Validator(validator httpx.RequestValidator) httpx.Middleware {
	errResponse := httpx.ServerErrorResponses(
		httpx.MakeServerError,
		httpx.NewErrorToStatusMapper(httpx.DefaultErrorToStatusMapping()),
	)

	return httpx.MiddlewareFn(
		func(response http.ResponseWriter, request *http.Request, next http.Handler) {
			ctx := request.Context()
			logger := logs.FromContext(ctx)

			if err := validator.Validate(request); err != nil {
				if httpx.IsParseError(err) {
					httpx.SendResponse(
						ctx, response, errResponse(errs.Wrapf(errs.ErrIllegalArgument, "failed to parse request body")),
					)

					return
				}

				if validatorErr, ok := httpx.ValidatorError(err); ok {
					logger.Warn().Err(errors.Unwrap(err)).Msgf("%s validation failed", validatorErr.Parameter)

					switch {
					case validatorErr.Parameter.IsQuery():
						err = errs.Wrapf(
							errs.ErrIllegalArgument, "incorrect url query \"%s\": %v", validatorErr.Field, validatorErr,
						)
					case validatorErr.Parameter.IsPath():
						err = errs.Wrapf(errs.ErrNotFound, "requested resource not found")
					case validatorErr.Field == "":
						err = errs.Wrapf(errs.ErrIllegalArgument, validatorErr.Message) // nolint:errorlint
					case validatorErr.Field != "":
						err = errs.WrapFields(
							errs.ErrIllegalArgument, validatorErr.Message, validatorErr.Field,
						) // nolint:errorlint
					}

					httpx.SendResponse(ctx, response, errResponse(err))

					return
				}

				logger.Err(err).Msg("got error while performing request validation with middleware")

				httpx.SendResponse(ctx, response, errResponse(err))

				return
			}

			next.ServeHTTP(response, request)
		},
	).Middleware()
}

func SecureValidator(validator httpx.RequestValidator) httpx.Middleware {
	errResponse := httpx.ServerErrorResponses(
		httpx.MakeSecureValidatorServerError,
		httpx.NewErrorToStatusMapper(httpx.DefaultErrorToStatusMapping()),
	)

	return httpx.MiddlewareFn(
		func(response http.ResponseWriter, request *http.Request, next http.Handler) {
			ctx := request.Context()
			logger := logs.FromContext(ctx)

			if err := validator.Validate(request); err != nil {
				if httpx.IsParseError(err) {
					httpx.SendResponse(
						ctx, response, errResponse(errs.Wrapf(errs.ErrIllegalArgument, "failed to parse request body")),
					)

					return
				}

				if validatorErr, ok := httpx.ValidatorError(err); ok {
					logger.Warn().Err(errors.Unwrap(err)).Msgf("%s validation failed", validatorErr.Parameter)

					switch {
					case validatorErr.Parameter.IsQuery():
						err = errs.Wrapf(
							errs.ErrIllegalArgument, "incorrect url query \"%s\": %v", validatorErr.Field, validatorErr,
						)
					case validatorErr.Parameter.IsPath():
						err = errs.Wrapf(errs.ErrNotFound, "requested resource not found")
					case validatorErr.Field == "":
						err = errs.Wrapf(errs.ErrIllegalArgument, validatorErr.Message) // nolint:errorlint
					case validatorErr.Field != "":
						err = errs.WrapFields(
							errs.ErrIllegalArgument, validatorErr.Message, validatorErr.Field,
						) // nolint:errorlint
					}

					httpx.SendResponse(ctx, response, errResponse(err))

					return
				}

				logger.Err(err).Msg("got error while performing request validation with middleware")

				httpx.SendResponse(ctx, response, errResponse(err))

				return
			}

			next.ServeHTTP(response, request)
		},
	).Middleware()
}
