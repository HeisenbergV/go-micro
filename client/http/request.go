package http

import "core/client"

type httpRequest struct {
	service     string
	method      string
	contentType string
	request     interface{}
}

func newHTTPRequest(service, method string, request interface{}, contentType string) client.Request {
	return &httpRequest{
		service:     service,
		method:      method,
		request:     request,
		contentType: contentType,
	}
}

func (h *httpRequest) ContentType() string {
	return h.contentType
}

func (h *httpRequest) Service() string {
	return h.service
}

func (h *httpRequest) Method() string {
	return h.method
}

func (h *httpRequest) Endpoint() string {
	return h.method
}

func (h *httpRequest) Body() interface{} {
	return h.request
}
