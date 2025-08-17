// Package httpx расширяет стандартный пакет net/http функциями для работы с протоколом.
//
// Пакет содержит хелперы, призванные сократить бойлерплейт код в приложении при работе с http: чтение запросов, запись
// ответов, работа с заголовками, посредники (middlewares).
package httpx

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/wal1251/pkg/core/presenters"
	"github.com/wal1251/pkg/core/security"
)

const (
	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"

	ContentTypeJSON = "application/json"
	ContentTypeXML  = "application/xml"
	BearerKeyword   = "Bearer"
)

var (
	_ presenters.StringViewer                                = Header{}
	_ presenters.StringViewer                                = (*Request)(nil)
	_ security.HTTPCredentialsProvider[security.BearerToken] = BearerTokenExtract
)

type (
	Header  http.Header
	Request http.Request
)

const ResponseStatusDefault = http.StatusOK

func (h Header) ContentType() string {
	return h.AsHTTP().Get(HeaderContentType)
}

func (h Header) AsHTTP() http.Header {
	return (http.Header)(h)
}

func (h Header) HasJSONContent() bool {
	return strings.Contains(h.AsHTTP().Get(HeaderContentType), ContentTypeJSON)
}

func (h Header) HasXMLContent() bool {
	return strings.Contains(h.AsHTTP().Get(HeaderContentType), ContentTypeXML)
}

func (h Header) StringView(view presenters.ViewType, opts presenters.ViewOptions) string {
	pairs := make(map[string]string)

	for key, list := range h {
		s := strings.Join(list, "|")
		if view == presenters.ViewLogs && key == HeaderAuthorization {
			s = "*"
		}

		pairs[key] = s
	}

	return presenters.ParameterView(pairs, view, opts)
}

func (h Header) InterfaceView(view presenters.ViewType, _ presenters.ViewOptions) any {
	pairs := make(map[string]string)

	for key, list := range h {
		s := strings.Join(list, "|")
		if view == presenters.ViewLogs && key == HeaderAuthorization {
			s = "*"
		}

		pairs[key] = s
	}

	return pairs
}

func (r *Request) ReadBody() ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read request body: %w", err)
	}

	defer func() {
		_ = r.Body.Close()
	}()

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return body, nil
}

func (r *Request) StringView(_ presenters.ViewType, _ presenters.ViewOptions) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf("Получен запрос %s %s://%s%s",
		r.Method, scheme, r.Host, r.RequestURI)
}

func BearerTokenExtract(r *http.Request) security.BearerToken {
	token := r.Header.Get(HeaderAuthorization)
	if len(token) > 0 {
		parts := strings.Split(token, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], BearerKeyword) {
			return ""
		}

		return security.BearerToken(parts[1])
	}

	return ""
}

func BearerTokenSet(r *http.Request, token security.BearerToken) {
	r.Header.Set(HeaderAuthorization, fmt.Sprintf("%s %s", BearerKeyword, token))
}
