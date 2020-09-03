package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	gmetadata "google.golang.org/grpc/metadata"
)

type grpcClient struct {
	opts Options
	pool *Pool
}

func (g *grpcClient) call(ctx context.Context, req Request, rsp interface{}, opts CallOptions) error {
	var header map[string]string

	header = make(map[string]string)

	header["timeout"] = fmt.Sprintf("%d", opts.RequestTimeout)
	header["x-content-type"] = req.ContentType()

	md := gmetadata.New(header)
	ctx = gmetadata.NewOutgoingContext(ctx, md)

	cf, err := g.newGRPCCodec(req.ContentType())
	if err != nil {
		return errors.New("[grpc] codec error:" + err.Error())
	}

	cc, err := g.pool.Get(ctx)
	if err != nil {
		return errors.New("[grpc] get pool rpc client error: " + err.Error())
	}
	defer func() {
		cc.Close()
	}()

	ch := make(chan error, 1)

	go func() {
		grpcCallOptions := []grpc.CallOption{
			grpc.ForceCodec(cf),
			grpc.CallContentSubtype(cf.Name())}

		err := cc.Invoke(ctx, methodToGRPC(req.Service(), req.Endpoint()), req.Body(), rsp, grpcCallOptions...)
		ch <- err
	}()

	var grr error
	select {
	case err := <-ch:
		grr = err
	case <-ctx.Done():
		grr = errors.New("[grpc] grpc reqeust timeout")
	}

	return grr
}

func (g *grpcClient) newGRPCCodec(contentType string) (encoding.Codec, error) {
	codecs := make(map[string]encoding.Codec)

	if c, ok := codecs[contentType]; ok {
		return wrapCodec{c}, nil
	}
	if c, ok := defaultGRPCCodecs[contentType]; ok {
		return wrapCodec{c}, nil
	}
	return nil, fmt.Errorf("Unsupported Content-Type: %s", contentType)
}

func (g *grpcClient) Init(opts ...Option) error {
	for _, o := range opts {
		o(&g.opts)
	}

	return nil
}

func (r *grpcClient) Options() Options {
	return r.opts
}

func (g *grpcClient) NewRequest(service, method string, req interface{}, reqOpts ...RequestOption) Request {
	return newGRPCRequest(service, method, req, g.opts.ContentType)
}

func (g *grpcClient) Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error {
	if req == nil {
		return errors.New("[grpc] request is nil")
	} else if rsp == nil {
		return errors.New("[grpc] response is nil")
	}

	callOpts := g.opts.CallOptions
	for _, opt := range opts {
		opt(&callOpts)
	}

	d, ok := ctx.Deadline()
	if !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, callOpts.RequestTimeout)
		defer cancel()
	} else {
		callOpts.RequestTimeout = time.Until(d)
	}

	select {
	case <-ctx.Done():
		return errors.New("[grpc] request timeout - deadline")
	default:
	}

	gcall := g.call

	//包装请求
	for i := len(callOpts.CallWrappers); i > 0; i-- {
		gcall = callOpts.CallWrappers[i-1](gcall)
	}

	call := func(i int) error {
		err := gcall(ctx, req, rsp, callOpts)
		return err
	}

	ch := make(chan error, callOpts.Retries+1)
	var gerr error

	for i := 0; i <= callOpts.Retries; i++ {
		go func(i int) {
			ch <- call(i)
		}(i)

		select {
		case <-ctx.Done():
			return errors.New("[grpc] request timeout")
		case err := <-ch:
			if err == nil {
				return nil
			}
			gerr = err
		}
	}

	return gerr
}

func (r *grpcClient) Publish(ctx context.Context, msg Message, opts ...PublishOption) error {
	return nil
}

func (g *grpcClient) String() string {
	return "grpc"
}

func newRpcClient(opts ...Option) Client {
	opt := NewOptions(opts...)

	pool, err := NewPool(
		func() (*grpc.ClientConn, error) {
			return grpc.Dial(opt.Address, grpc.WithInsecure())
		}, opt.PoolInitNum, opt.PoolCapacity, opt.PoolTTL, 0)

	if err != nil {
		return nil
	}

	rc := &grpcClient{
		opts: opt,
	}

	rc.pool = pool
	c := Client(rc)

	for i := len(opt.Wrappers); i > 0; i-- {
		c = opt.Wrappers[i-1](c)
	}

	return c
}

func NewRpcClient(opts ...Option) Client {
	return newRpcClient(opts...)
}
