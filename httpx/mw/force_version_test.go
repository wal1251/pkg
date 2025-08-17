package mw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func TestMinVersion(t *testing.T) {
	t.Run("good case", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set(HeaderFrontVersion, "100")

		fHandle := MinVersion("100")
		handle := fHandle(handler{})

		handle.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", w.Body.String())
	})
	t.Run("good case semver1", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set(HeaderFrontVersion, "1.0.0")

		fHandle := MinVersion("1.0.0")
		handle := fHandle(handler{})

		handle.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", w.Body.String())
	})
	t.Run("good case semver2", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set(HeaderFrontVersion, "1.0.10")

		fHandle := MinVersion("1.0.0.1")
		handle := fHandle(handler{})

		handle.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", w.Body.String())
	})
	t.Run("good case hardcoded", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set(HeaderFrontVersion, "0.0.0.0")

		fHandle := MinVersion("1.0.1")
		handle := fHandle(handler{})

		handle.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", w.Body.String())
	})
	t.Run("good case with 0 version & no header", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/", nil)

		fHandle := MinVersion("0")
		handle := fHandle(handler{})

		handle.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", w.Body.String())
	})
	t.Run("no version in header", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/", nil)

		fHandle := MinVersion("100")
		handle := fHandle(handler{})

		handle.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"error": "1.2"}`, w.Body.String())
	})
	t.Run("hardcoded error", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set(HeaderFrontVersion, "0.0.0")

		fHandle := MinVersion("1.0.0")
		handle := fHandle(handler{})

		handle.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"error": "1.2"}`, w.Body.String())
	})
	t.Run("low version in header", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set(HeaderFrontVersion, "99")

		fHandle := MinVersion("100")
		handle := fHandle(handler{})

		handle.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"error": "1.2"}`, w.Body.String())
	})
	t.Run("no header", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/", nil)

		fHandle := MinVersion("100")
		handle := fHandle(handler{})

		handle.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"error": "1.2"}`, w.Body.String())
	})
}
