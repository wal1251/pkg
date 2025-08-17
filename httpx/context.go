package httpx

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

const (
	xForwardedFor  = "X-Forwarded-For"
	xForwardedHost = "X-Forwarded-Host"
	xRequestID     = "X-Request-ID"
)

// AnnotateContext добавляет необходимые заголовки запроса в метадату grpc контекста
// по аналогии метода из grpc-gateway/runtime/context.go.
func AnnotateContext(req *http.Request) context.Context {
	var pairs []string
	if host := req.Header.Get(xForwardedHost); host != "" {
		pairs = append(pairs, strings.ToLower(xForwardedHost), host)
	} else if req.Host != "" {
		pairs = append(pairs, strings.ToLower(xForwardedHost), req.Host)
	}

	xff := req.Header.Values(xForwardedFor)
	if addr := req.RemoteAddr; addr != "" {
		if remoteIP, _, err := net.SplitHostPort(addr); err == nil {
			xff = append(xff, remoteIP)
		}
	}
	if len(xff) > 0 {
		pairs = append(pairs, strings.ToLower(xForwardedFor), strings.Join(xff, ", "))
	}

	reqID := req.Header.Get(xRequestID)
	if reqID == "" {
		reqID = uuid.New().String()
	}
	pairs = append(pairs, strings.ToLower(xRequestID), reqID)

	return metadata.NewOutgoingContext(req.Context(), metadata.Pairs(pairs...))
}
