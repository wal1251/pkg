// Package mw содержит готовые реализации наиболее часто используемых http посредников (middlewares).
package mw

import (
	"net/http"

	"github.com/wal1251/pkg/httpx"
)

func Middleware(mw httpx.Middleware) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mw(http.HandlerFunc(next.ServeHTTP)).ServeHTTP(w, r)
		})
	}
}
