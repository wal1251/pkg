package mw_test

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/httpx"
	"github.com/wal1251/pkg/httpx/mw"
	"github.com/wal1251/pkg/tools/collections"
	"github.com/wal1251/pkg/tools/crypto"
)

func TestRequestSignature_ValidSignature(t *testing.T) {
	secret := "testsecret"
	requestLifetime := time.Hour
	skipURLs := collections.NewSet[string]()

	endpointHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.SendResponse(r.Context(), w,
			httpx.NewServerResponse[map[string]string]().
				WithContentTypeJSON().
				WithValue(map[string]string{"status": "success"}))
	})

	// Генерация валидной подписи
	method := http.MethodGet
	url := "/"
	creationTime := time.Now().Unix()
	creationTimeEncoded := base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(creationTime, 10)))
	hmac := crypto.NewHMAC(secret)
	signatureBase := fmt.Sprintf("%s%s%s", method, url, strconv.FormatInt(creationTime, 10))
	signatureHash := base64.StdEncoding.EncodeToString([]byte(hmac.Sign(signatureBase)))

	r := httptest.NewRequest(method, url, nil)
	r.Header.Set("Request-Signature", signatureHash)
	r.Header.Set("Request-Creation-Time", creationTimeEncoded)

	wr := httptest.NewRecorder()

	rs := mw.RequestSignature(secret, requestLifetime, skipURLs)
	rs(endpointHandler).ServeHTTP(wr, r)

	assert.Equal(t, http.StatusOK, wr.Code, "Expected HTTP status OK")

	expectedResponseBody := `{"status":"success"}`
	assert.Equal(t, expectedResponseBody, strings.TrimSpace(wr.Body.String()), "Expected body to match exactly")
}

func TestRequestSignature_InvalidSignature(t *testing.T) {
	secret := "testsecret"
	requestLifetime := time.Hour
	skipURLs := collections.NewSet[string]()

	endpointHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("Endpoint handler should not be called")
	})

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Request-Signature", "invalid_signature")
	r.Header.Set("Request-Creation-Time", base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(time.Now().Unix(), 10))))

	wr := httptest.NewRecorder()

	rs := mw.RequestSignature(secret, requestLifetime, skipURLs)
	rs(endpointHandler).ServeHTTP(wr, r)

	assert.Equal(t, http.StatusForbidden, wr.Code, "Expected HTTP status forbidden")

	expectedResponseBody := `{"code":"FORBIDDEN","message":"FORBIDDEN: failed to parse request signature"}`
	assert.Equal(t, expectedResponseBody, strings.TrimSpace(wr.Body.String()), "Expected body to match error message")
}

func TestRequestSignature_RequestExpired(t *testing.T) {
	secret := "testsecret"
	requestLifetime := 1 * time.Second
	skipURLs := collections.NewSet[string]()

	endpointHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("Endpoint handler should not be called for expired request")
	})

	expiredTime := time.Now().Add(-2 * requestLifetime).Unix()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Request-Signature", "123") // Просто заглушка, т.к. проверка времени идёт раньше
	r.Header.Set("Request-Creation-Time", base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(expiredTime, 10))))

	wr := httptest.NewRecorder()

	rs := mw.RequestSignature(secret, requestLifetime, skipURLs)
	rs(endpointHandler).ServeHTTP(wr, r)

	assert.Equal(t, http.StatusForbidden, wr.Code, "Expected HTTP status forbidden")

	expectedResponseBody := `{"code":"FORBIDDEN","message":"FORBIDDEN: request is outdated"}`
	assert.Equal(t, expectedResponseBody, strings.TrimSpace(wr.Body.String()), "Expected body to match error message for expired request")
}

func TestRequestSignature_SkipURLs(t *testing.T) {
	secret := "testsecret"
	requestLifetime := time.Hour
	skippedPath := "/skip"
	skipURLs := collections.NewSet[string](skippedPath)

	endpointHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.SendResponse(r.Context(), w,
			httpx.NewServerResponse[map[string]string]().
				WithContentTypeJSON().
				WithValue(map[string]string{"skipped": "true"}))
	})

	r := httptest.NewRequest(http.MethodGet, skippedPath, nil)

	wr := httptest.NewRecorder()

	rs := mw.RequestSignature(secret, requestLifetime, skipURLs)
	rs(endpointHandler).ServeHTTP(wr, r)

	assert.Equal(t, http.StatusOK, wr.Code, "Expected HTTP status OK for skipped URL")

	expectedResponseBody := `{"skipped":"true"}`
	assert.Equal(t, expectedResponseBody, strings.TrimSpace(wr.Body.String()), "Expected body to match exactly")
}

func TestRequestSignature_DecodingCreationTimeError(t *testing.T) {
	secret := "testsecret"
	requestLifetime := time.Hour
	skipURLs := collections.NewSet[string]()

	endpointHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("Endpoint handler should not be called when there's an error decoding creation time")
	})

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Request-Signature", "123") // Просто заглушка, т.к. ошибка будет на этапе декодирования времени
	r.Header.Set("Request-Creation-Time", "invalid_base64")

	wr := httptest.NewRecorder()

	rs := mw.RequestSignature(secret, requestLifetime, skipURLs)
	rs(endpointHandler).ServeHTTP(wr, r)

	assert.Equal(t, http.StatusForbidden, wr.Code, "Expected HTTP status forbidden")

	expectedResponseBody := `{"code":"FORBIDDEN","message":"FORBIDDEN: failed to parse request time of creation"}`
	assert.Equal(t, expectedResponseBody, strings.TrimSpace(wr.Body.String()), "Expected body to match error message for decoding error")
}
