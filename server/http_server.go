package server

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

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
	sync.Mutex
	opts         Options
	hd           Handler
	exit         chan chan error
	registerOnce sync.Once
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
	h.Lock()
	opts := h.opts
	h.Unlock()
	return opts
}

func (h *httpServer) Init(opts ...Option) error {
	h.Lock()
	for _, o := range opts {
		o(&h.opts)
	}
	h.Unlock()
	return nil
}

func (h *httpServer) Handle(handler Handler) error {
	if _, ok := handler.Handler().(http.Handler); !ok {
		return errors.New("Handle requires http.Handler")
	}
	h.Lock()
	h.hd = handler
	h.Unlock()
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
	h.Lock()
	opts := h.opts
	hd := h.hd
	h.Unlock()

	ln, err := net.Listen("tcp", opts.Address)
	if err != nil {
		return err
	}

	h.Lock()
	h.opts.Address = ln.Addr().String()
	h.Unlock()

	handler, ok := hd.Handler().(http.Handler)
	if !ok {
		return errors.New("Server required http.Handler")
	}

	go http.Serve(ln, handler)

	return nil
}

func (h *httpServer) Stop() error {
	ch := make(chan error)
	h.exit <- ch
	return <-ch
}

func (h *httpServer) String() string {
	return "http"
}

func newHttpServer(opts ...Option) Server {
	return &httpServer{
		opts: newOptions(opts...),
		exit: make(chan chan error),
	}
}
