package mw

import (
	"context"
	"net/http"

	"github.com/wal1251/pkg/httpx"
	"github.com/wal1251/pkg/tools/acceptlanguage"
)

const (
	HeaderAcceptLanguage = "Accept-Language"
)

func AcceptLanguage() httpx.Middleware {
	return httpx.MiddlewareFn(func(response http.ResponseWriter, request *http.Request, next http.Handler) {
		acceptLanguage := acceptlanguage.Validate(request.Header.Get(HeaderAcceptLanguage))
		ctx := context.WithValue(request.Context(), acceptlanguage.AcceptLanguageKey, acceptLanguage)
		next.ServeHTTP(response, request.WithContext(ctx))
	}).Middleware()
}
