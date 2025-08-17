package httpx_test

import (
	"bytes"
	"encoding/xml"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/httpx"
)

func TestServerRequest_Decode_json(t *testing.T) {
	type Sample struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	r := httptest.NewRequest("POST", "/foo", bytes.NewBufferString(`{"foo":"baz", "bar": 1}`))
	r.Header.Add(httpx.HeaderContentType, "application/json;charset=UTF-8")

	request := httpx.NewServerRequest[Sample](r)
	if assert.NoError(t, request.Decode()) {
		assert.Equal(t, Sample{Foo: "baz", Bar: 1}, request.Value)
	}
}

func TestServerRequest_Decode_xml(t *testing.T) {
	type Sample struct {
		XMLName xml.Name `xml:"Sample"`
		Foo     string   `xml:"Foo"`
		Bar     int      `xml:"Bar"`
	}

	r := httptest.NewRequest("POST", "/foo", bytes.NewBufferString("<Sample><Foo>Hello</Foo><Bar>42</Bar></Sample>"))
	r.Header.Add(httpx.HeaderContentType, "application/xml;charset=UTF-8")

	request := httpx.NewServerRequest[Sample](r)
	if assert.NoError(t, request.Decode()) {
		assert.Equal(t, Sample{XMLName: xml.Name{Space: "", Local: "Sample"}, Foo: "Hello", Bar: 42}, request.Value)
	}
}

func TestServerRequest_Decode_void(t *testing.T) {
	r := httptest.NewRequest("GET", "/foo", nil)
	request := httpx.NewServerRequest[any](r)
	if assert.NoError(t, request.Decode()) {
		assert.Nil(t, request.Value)
	}
}
