package mw

import (
	"errors"
	"net/http"

	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/httpx"
)

func Recoverer() httpx.Middleware {
	return httpx.MiddlewareFn(func(response http.ResponseWriter, request *http.Request, next http.Handler) {
		defer func() {
			if artifact := recover(); artifact != nil {
				if err, ok := artifact.(error); ok {
					if errors.Is(err, http.ErrAbortHandler) {
						return
					}
				}

				logs.FromContext(request.Context()).
					Error().Stack().Msgf("panic: %v", artifact)

				// Не пишем тело ответа, т.к. в w уже могли что-то записать ранее.
				response.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(response, request)
	}).Middleware()
}
