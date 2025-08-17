package mw_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wal1251/pkg/httpx/mw"
	"github.com/wal1251/pkg/tools/acceptlanguage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcceptLanguage(t *testing.T) {
	t.Run("checking the completed header", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)
		req.Header = http.Header{
			"Accept-Language": []string{"kk"},
		}

		rr := httptest.NewRecorder()

		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acceptLanguage := r.Context().Value(acceptlanguage.AcceptLanguageKey)
			require.Equal(t, "kk", acceptLanguage, "AcceptLanguage should be present in context")
			w.WriteHeader(http.StatusOK)
		})

		mw.AcceptLanguage()(mockHandler).ServeHTTP(rr, req)
		// Проверка на 200 OK
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("checking empty header", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acceptLanguage := r.Context().Value(acceptlanguage.AcceptLanguageKey)
			require.Equal(t, "ru", acceptLanguage, "AcceptLanguage should be present in context")
			w.WriteHeader(http.StatusOK)
		})

		mw.AcceptLanguage()(mockHandler).ServeHTTP(rr, req)
		// Проверка на 200 OK
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
