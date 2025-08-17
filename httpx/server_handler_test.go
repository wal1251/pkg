package httpx_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/core/errs"
	"github.com/wal1251/pkg/httpx"
)

func TestServerHandler_Handle_json(t *testing.T) {
	type SampleRequest struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	type SampleResponse struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/foo", bytes.NewBufferString(`{"foo":"baz", "bar": 1}`))
	r.Header.Add(httpx.HeaderContentType, "application/json;charset=UTF-8")

	httpx.NewServerHandler[SampleRequest, SampleResponse](w, r).
		WithResponseJSON().
		WithResponseStatus(http.StatusAccepted).
		WithMethod(func(ctx context.Context, req SampleRequest) (SampleResponse, error) {
			require.Equalf(t, "baz", req.Foo, "request argument doesn't match expected")
			require.Equalf(t, 1, req.Bar, "request argument doesn't match expected")

			return SampleResponse{
				Foo: req.Foo + "1",
				Bar: req.Bar + 1,
			}, nil
		}).
		Handle()

	if assert.Equalf(t, http.StatusAccepted, w.Code, "response status code not matches") {
		assert.JSONEqf(t, `{"foo":"baz1", "bar": 2}`, w.Body.String(), "response body not matches")
	}
}

func TestServerHandler_Handle_xml(t *testing.T) {
	type SampleRequest struct {
		XMLName xml.Name `xml:"SampleRequest"`
		Foo     string   `xml:"Foo"`
		Bar     int      `xml:"Bar"`
	}

	type SampleResponse struct {
		XMLName xml.Name `xml:"SampleResponse"`
		Foo     string   `xml:"Foo"`
		Bar     int      `xml:"Bar"`
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/foo", bytes.NewBufferString("<SampleRequest><Foo>Hello</Foo><Bar>42</Bar></SampleRequest>"))
	r.Header.Set(httpx.HeaderContentType, "application/xml;charset=UTF-8")

	httpx.NewServerHandler[SampleRequest, SampleResponse](w, r).
		WithResponseXML().
		WithResponseStatus(http.StatusAccepted).
		WithMethod(func(ctx context.Context, req SampleRequest) (SampleResponse, error) {
			require.Equal(t, "Hello", req.Foo, "request argument 'Foo' doesn't match expected")
			require.Equal(t, 42, req.Bar, "request argument 'Bar' doesn't match expected")

			return SampleResponse{
				Foo: req.Foo + "1",
				Bar: req.Bar + 1,
			}, nil
		}).
		Handle()

	if assert.Equal(t, http.StatusAccepted, w.Code, "response status code does not match") {
		assert.Equal(
			t,
			"<SampleResponse><Foo>Hello1</Foo><Bar>43</Bar></SampleResponse>", w.Body.String(),
			"response body does not match")
	}
}

func TestServerHandler_Handle_error(t *testing.T) {
	type Sample struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	errMap := httpx.NewErrorToStatusMapper(httpx.DefaultErrorToStatusMapping())

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/foo", bytes.NewBufferString(`{"foo":"baz", "bar": 1}`))
	r.Header.Add(httpx.HeaderContentType, "application/json;charset=UTF-8")

	httpx.NewServerHandler[Sample, any](w, r).
		WithMethod(func(ctx context.Context, req Sample) (any, error) {
			return nil, errs.Wrapf(errs.ErrIllegalArgument, "fake error")
		}).
		Handle()

	if assert.Equalf(t, errMap.Status(errs.ErrIllegalArgument.Type), w.Code, "response status code not matches") {
		assert.JSONEqf(t, `{"code":"ILLEGAL_ARGUMENT", "message":"ILLEGAL_ARGUMENT: fake error"}`, w.Body.String(), "response body not matches")
	}
}

func TestServerHandler_Handle_void(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/foo", nil)

	httpx.NewServerHandler[any, any](w, r).
		WithMethod(func(ctx context.Context, req any) (any, error) {
			require.Nilf(t, req, "request argument doesn't match expected")
			return nil, nil
		}).
		Handle()

	if assert.Equalf(t, httpx.ResponseStatusDefault, w.Code, "response status code not matches") {
		assert.Equalf(t, 0, w.Body.Len(), "response body not matches")
	}
}
