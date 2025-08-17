package httpx_test

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/core/errs"
	"github.com/wal1251/pkg/httpx"
)

func TestServerResponse_json(t *testing.T) {
	type Sample struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	w := httptest.NewRecorder()

	err := httpx.NewServerResponse[Sample]().
		WithContentTypeJSON().
		WithStatus(http.StatusCreated).
		WithValue(Sample{Foo: "bar", Bar: 99}).
		Send(w)

	if assert.NoError(t, err) {
		assert.Equalf(t, http.StatusCreated, w.Code, "response status code not matches")
		assert.JSONEqf(t, `{"foo":"bar", "bar": 99}`, w.Body.String(), "response body not matches")
	}
}

func TestServerResponse_xml(t *testing.T) {
	type Sample struct {
		XMLName xml.Name `xml:"Sample"`
		Foo     string   `xml:"Foo"`
		Bar     int      `xml:"Bar"`
	}

	w := httptest.NewRecorder()
	err := httpx.NewServerResponse[Sample]().
		WithContentTypeXML().
		WithStatus(http.StatusCreated).
		WithValue(Sample{Foo: "Hello", Bar: 42}).
		Send(w)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, w.Code, "response status code not matches")
		assert.Equal(t, "<Sample><Foo>Hello</Foo><Bar>42</Bar></Sample>", w.Body.String(), "response body not matches")
	}
}

func TestServerResponse_empty_body(t *testing.T) {
	w := httptest.NewRecorder()

	err := httpx.NewServerResponse[any]().
		WithStatus(http.StatusNoContent).
		Send(w)

	if assert.NoError(t, err) {
		assert.Equalf(t, http.StatusNoContent, w.Code, "response status code not matches")
		assert.Equalf(t, 0, w.Body.Len(), "response body not matches")
	}
}

func TestErrorServerResponses(t *testing.T) {
	errMapper := httpx.NewErrorToStatusMapper(httpx.DefaultErrorToStatusMapping())

	makeErrorResponse := httpx.ServerErrorResponses(
		httpx.MakeServerError,
		httpx.NewErrorToStatusMapper(httpx.DefaultErrorToStatusMapping()),
	)

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Unclassified server error",
			err:        errors.New("fake error"),
			wantStatus: errMapper.Status(errs.ErrSystemFailure.Type),
			wantBody:   `{"code":"SYSTEM_FAILURE", "message":"fake error"}`,
		},
		{
			name:       "Described server error",
			err:        errs.Wrapf(errs.ErrIllegalArgument, "fake error"),
			wantStatus: errMapper.Status(errs.ErrIllegalArgument.Type),
			wantBody:   `{"code":"ILLEGAL_ARGUMENT", "message":"ILLEGAL_ARGUMENT: fake error"}`,
		},
		{
			name:       "Void server error",
			err:        errs.Error{},
			wantStatus: errMapper.Status(errs.Error{}.Type),
			wantBody:   `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if assert.NoError(t, makeErrorResponse(tt.err).Send(w)) {
				assert.Equalf(t, tt.wantStatus, w.Code, "response status code not matches")
				assert.JSONEqf(t, tt.wantBody, w.Body.String(), "response body not matches")
			}
		})
	}
}
