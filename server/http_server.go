package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"core/codec"
	"core/codec/jsonrpc"
	"core/codec/protorpc"
)

var (
	defaultCodecs = map[string]codec.NewCodec{
		"application/json":         jsonrpc.NewCodec,
		"application/json-rpc":     jsonrpc.NewCodec,
		"application/protobuf":     protorpc.NewCodec,
		"application/proto-rpc":    protorpc.NewCodec,
		"application/octet-stream": protorpc.NewCodec,
	}
)

type httpServer struct {
	opts Options
	http *http.Server
}

func (h *httpServer) newCodec(contentType string) (codec.NewCodec, error) {
	if cf, ok := h.opts.Codecs[contentType]; ok {
		return cf, nil
	}
	if cf, ok := defaultCodecs[contentType]; ok {
		return cf, nil
	}
	return nil, fmt.Errorf("Unsupported Content-Type: %s", contentType)
}

func (h *httpServer) Options() Options {
	opts := h.opts
	return opts
}

func (h *httpServer) Init(opts ...Option) error {
	for _, o := range opts {
		o(&h.opts)
	}

	h.http = &http.Server{}
	h.http.Addr = h.opts.Address
	return nil
}

func (h *httpServer) Handle(handler Handler) error {
	if _, ok := handler.Handler().(http.Handler); !ok {
		return errors.New("Handle requires http.Handler")
	}

	h.http.Handler = handler.Handler().(http.Handler)
	return nil
}

func (h *httpServer) NewHandler(handler interface{}, opts ...HandlerOption) Handler {
	options := HandlerOptions{
		Metadata: make(map[string]map[string]string),
	}

	for _, o := range opts {
		o(&options)
	}

	return &httpHandler{
		hd:   handler,
		opts: options,
	}
}

func (h *httpServer) Start() error {

	go h.http.ListenAndServe()

	return nil
}

func (h *httpServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return h.http.Shutdown(ctx)
}

func (h *httpServer) String() string {
	return "http"
}

func newHttpServer(opts ...Option) Server {
	return &httpServer{
		opts: newOptions(opts...),
	}
}
