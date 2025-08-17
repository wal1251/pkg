package httpx

import (
	"net/http"

	"github.com/wal1251/pkg/core"
	"github.com/wal1251/pkg/core/errs"
	"github.com/wal1251/pkg/core/logs"
)

type (
	ServerHandler[REQ, RESP any] struct {
		w         http.ResponseWriter
		method    core.Method[REQ, RESP]
		request   *ServerRequest[REQ]
		response  *ServerResponseBuilder[RESP]
		errHandle func(err error) *ServerResponseBuilder[ServerError]
	}
)

func (h *ServerHandler[REQ, RESP]) WithMethod(method core.Method[REQ, RESP]) *ServerHandler[REQ, RESP] {
	h.method = method

	return h
}

func (h *ServerHandler[REQ, RESP]) WithResponse(opt func(*ServerResponseBuilder[RESP])) *ServerHandler[REQ, RESP] {
	opt(h.response)

	return h
}

func (h *ServerHandler[REQ, RESP]) WithErrHandler(errHandler func(err error) *ServerResponseBuilder[ServerError]) *ServerHandler[REQ, RESP] {
	h.errHandle = errHandler

	return h
}

func (h *ServerHandler[REQ, RESP]) WithResponseJSON() *ServerHandler[REQ, RESP] {
	h.response.WithContentTypeJSON()

	return h
}

func (h *ServerHandler[REQ, RESP]) WithResponseXML() *ServerHandler[REQ, RESP] {
	h.response.WithContentTypeXML()

	return h
}

func (h *ServerHandler[REQ, RESP]) WithResponseStatus(status int) *ServerHandler[REQ, RESP] {
	h.response.WithStatus(status)

	return h
}

func (h *ServerHandler[REQ, RESP]) Handle() {
	ctx := h.request.Context()
	log := logs.FromContext(ctx)

	if err := h.request.Decode(); err != nil {
		log.Err(err).Msg("can't decode request body")

		if err := h.errHandle(
			errs.Wrapf(
				errs.ErrIllegalArgument, "can't decode argument: %v", err,
			),
		).Send(h.w); err != nil {
			log.Err(err).Msg("can't send error response")
		}

		return
	}

	result, err := h.method(ctx, h.request.Value)
	if err != nil {
		log.Err(err).Msg("error occurred while handling request")

		if err := h.errHandle(err).Send(h.w); err != nil {
			log.Err(err).Msg("can't send error response")
		}

		return
	}

	if err = h.response.WithValue(result).Send(h.w); err != nil {
		log.Err(err).Msg("can't send request body")
	}
}

func NewServerHandler[REQ, RESP any](writer http.ResponseWriter, request *http.Request) *ServerHandler[REQ, RESP] {
	return &ServerHandler[REQ, RESP]{
		w:         writer,
		request:   NewServerRequest[REQ](request),
		response:  NewServerResponse[RESP](),
		errHandle: ServerErrorResponses(MakeServerError, NewErrorToStatusMapper(DefaultErrorToStatusMapping())),
	}
}

func NewSecureServerHandler[REQ, RESP any](
	writer http.ResponseWriter,
	request *http.Request,
) *ServerHandler[REQ, RESP] {
	return &ServerHandler[REQ, RESP]{
		w:         writer,
		request:   NewServerRequest[REQ](request),
		response:  NewServerResponse[RESP](),
		errHandle: ServerErrorResponses(MakeSecureServerError, NewErrorToStatusMapper(DefaultErrorToStatusMapping())),
	}
}
