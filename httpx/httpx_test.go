package httpx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/core/security"
	"github.com/wal1251/pkg/httpx"
)

func TestBearerTokenAuthorizer(t *testing.T) {
	requestWithAuthorization := func(a string) *http.Request {
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Authorization", a)
		return r
	}

	tests := []struct {
		name string
		r    *http.Request
		want security.BearerToken
	}{
		{
			name: "Basic test",
			r:    requestWithAuthorization("Bearer FOO"),
			want: "FOO",
		},
		{
			name: "Request has no authorization",
			r:    httptest.NewRequest("GET", "/", nil),
			want: "",
		},
		{
			name: "Has no prefix word",
			r:    requestWithAuthorization("FOO"),
			want: "",
		},
		{
			name: "Has no delimiter space with prefix word",
			r:    requestWithAuthorization("BearerFOO"),
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, httpx.BearerTokenExtract(tt.r))
		})
	}
}
