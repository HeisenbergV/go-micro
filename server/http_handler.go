package server

type httpHandler struct {
	opts HandlerOptions
	hd   interface{}
}

func (h *httpHandler) Name() string {
	return "handler"
}

func (h *httpHandler) Handler() interface{} {
	return h.hd
}

func (h *httpHandler) Options() HandlerOptions {
	return h.opts
}
