package mw

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/wal1251/pkg/core/errs"
	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/httpx"
	"github.com/wal1251/pkg/tools/collections"
	"github.com/wal1251/pkg/tools/crypto"
)

// RequestSignature проверяет подпись каждого запроса с помощью hash функции.
func RequestSignature(secret string, requestLifetime time.Duration, skipURLs collections.Set[string]) httpx.Middleware {
	errResponse := httpx.ServerErrorResponses(
		httpx.MakeServerError,
		httpx.NewErrorToStatusMapper(httpx.DefaultErrorToStatusMapping()),
	)

	return httpx.MiddlewareFn(func(response http.ResponseWriter, request *http.Request, next http.Handler) {
		if skipURLs.Contains(request.URL.Path) {
			next.ServeHTTP(response, request)

			return
		}

		ctx := request.Context()
		logger := logs.FromContext(ctx)

		hmac := crypto.NewHMAC(secret)

		requestSignature := request.Header.Get("Request-Signature")
		requestCreationTime := request.Header.Get("Request-Creation-Time")

		creationTimeDecoded, err := base64.StdEncoding.DecodeString(requestCreationTime)
		if err != nil {
			logger.Error().Msg("failed to parse request time of creation")
			httpx.SendResponse(ctx, response, errResponse(errs.Wrapf(errs.ErrForbidden, "failed to parse request time of creation")))

			return
		}

		creationTime, err := strconv.ParseInt(string(creationTimeDecoded), 10, 64)
		if err != nil {
			logger.Error().Msg("failed to parse request time of creation")
			httpx.SendResponse(ctx, response, errResponse(errs.Wrapf(errs.ErrForbidden, "failed to parse request time of creation")))

			return
		}

		if time.Since(time.Unix(creationTime, 0)) > requestLifetime {
			logger.Error().Msg("request is outdated")
			httpx.SendResponse(ctx, response, errResponse(errs.Wrapf(errs.ErrForbidden, "request is outdated")))

			return
		}

		signatureBase := fmt.Sprintf("%s%s%s", request.Method, request.URL, string(creationTimeDecoded))
		signatureHash := base64.StdEncoding.EncodeToString([]byte(hmac.Sign(signatureBase)))

		if requestSignature != signatureHash {
			logger.Error().Msg("failed to parse request signature")
			httpx.SendResponse(ctx, response, errResponse(errs.Wrapf(errs.ErrForbidden, "failed to parse request signature")))

			return
		}

		next.ServeHTTP(response, request)
	}).Middleware()
}
