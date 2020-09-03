package http

import (
	"context"
	"errors"
	"fmt"
	"time"

	"core/client"

	"github.com/parnurzeal/gorequest"
)

type httpClient struct {
	opts client.Options
	hs   *gorequest.SuperAgent
}

func (h *httpClient) call(ctx context.Context, req client.Request, rsp interface{}, opts client.CallOptions) error {
	h.hs.Header.Add("timeout", fmt.Sprintf("%d", opts.RequestTimeout))
	h.hs.Header.Add("content-type", req.ContentType())

	var errs []error
	switch req.Method() {
	case "GET":
		_, _, errs = h.hs.Get(req.Service()).EndStruct(rsp)
		break
	case "POST":
		_, _, errs = h.hs.Post(req.Service()).Send(req).EndStruct(rsp)
		break
	case "PUT":
		_, _, errs = h.hs.Put(req.Service()).Send(req).EndStruct(rsp)
		break
	}

	if len(errs) != 0 {
		return errs[0]
	}
	return nil
}

func (h *httpClient) Init(opts ...client.Option) error {
	for _, o := range opts {
		o(&h.opts)
	}
	return nil
}

func (h *httpClient) NewRequest(service, method string, req interface{}, reqOpts ...client.RequestOption) client.Request {
	return newHTTPRequest(service, method, req, h.opts.ContentType)
}

func (h *httpClient) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	callOpts := h.opts.CallOptions
	for _, opt := range opts {
		opt(&callOpts)
	}

	d, ok := ctx.Deadline()
	if !ok {
		ctx, _ = context.WithTimeout(ctx, callOpts.RequestTimeout)
	} else {
		callOpts.RequestTimeout = time.Until(d)
	}

	select {
	case <-ctx.Done():
		return errors.New("[http] request timeout - deadline")
	default:
	}

	hcall := h.call

	for i := len(callOpts.CallWrappers); i > 0; i-- {
		hcall = callOpts.CallWrappers[i-1](hcall)
	}

	call := func(i int) error {
		err := hcall(ctx, req, rsp, callOpts)
		return err
	}

	ch := make(chan error, callOpts.Retries)
	var gerr error

	for i := 0; i < callOpts.Retries; i++ {
		go func() {
			ch <- call(i)
		}()

		select {
		case <-ctx.Done():
			return errors.New("[http] request timeout")
		case err := <-ch:
			if err == nil {
				return nil
			}
			gerr = err
		}
	}

	return gerr
}

func (h *httpClient) String() string {
	return "http"
}

func (r *httpClient) Options() client.Options {
	return r.opts
}

func (r *httpClient) Publish(ctx context.Context, msg client.Message, opts ...client.PublishOption) error {
	return nil
}

func newHttpClient(opts ...client.Option) client.Client {
	opt := client.NewOptions(opts...)

	hc := &httpClient{
		opts: opt,
		hs:   gorequest.New(),
	}

	c := client.Client(hc)

	for i := len(opt.Wrappers); i > 0; i-- {
		c = opt.Wrappers[i-1](c)
	}

	return c
}

func NewHttpClient(opts ...client.Option) client.Client {
	return newHttpClient(opts...)
}
