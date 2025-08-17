package httpx

import "net/http"

type (
	Middleware   func(next http.Handler) http.Handler
	MiddlewareFn func(w http.ResponseWriter, r *http.Request, next http.Handler)
)

func (f MiddlewareFn) Middleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			if f == nil {
				next.ServeHTTP(response, request)

				return
			}

			f(response, request, next)
		})
	}
}
