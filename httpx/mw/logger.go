package mw

import (
	"bytes"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/wal1251/pkg/core/ctxs"
	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/core/presenters"
	"github.com/wal1251/pkg/httpx"
)

func Logger(view presenters.ViewType, options presenters.ViewOptions) httpx.Middleware {
	return httpx.MiddlewareFn(func(response http.ResponseWriter, request *http.Request, next http.Handler) {
		ctx := ctxs.StartMeasureContext(request.Context())

		logger := logs.SubLogger(logs.FromContext(ctx))
		logger.UpdateContext(logs.Options(logs.WithRequestID(ctx),
			logs.FromTag.Option(request.RemoteAddr),
			logs.PathTag.Option(request.URL.Path),
			logs.MethodTag.Option(request.Method),
		))

		buf := bytes.NewBuffer([]byte{})
		responseWrapper := middleware.NewWrapResponseWriter(response, request.ProtoMajor)
		responseWrapper.Tee(buf)

		defer func() {
			logger.UpdateContext(logs.WithElapsedTime(ctx))
			logger.Info().
				Interface(string(logs.HeadersTag),
					httpx.Header(responseWrapper.Header()).InterfaceView(view, options)).
				Msgf("[%d]: %d bytes",
					responseWrapper.Status(),
					responseWrapper.BytesWritten(),
				)

			if responseWrapper.BytesWritten() > 0 {
				logger.Trace().Msgf("response: %s", presenters.JSONHideCredentials(buf.String(), options))
			}
		}()

		body, err := (*httpx.Request)(request).ReadBody()
		if err != nil {
			logger.Err(err).Msg("failed to read body")
			response.WriteHeader(http.StatusInternalServerError)

			return
		}

		logger.Info().
			Interface(string(logs.HeadersTag), httpx.Header(request.Header).InterfaceView(view, options)).
			Msgf(
				"%s: %d bytes",
				presenters.ParameterView((*httpx.Request)(request), view, options),
				len(body),
			)

		if len(body) != 0 && httpx.Header(request.Header).HasJSONContent() {
			logger.Trace().Msgf("request body: %s", presenters.JSONString(string(body), view, options))
		}

		next.ServeHTTP(responseWrapper, request.WithContext(logger.WithContext(ctx)))
	}).Middleware()
}
