package client

import (
	"context"
	"time"

	"core/broker"
	"core/codec"
)

type Options struct {
	Address      string
	ContentType  string
	UsePool      bool
	PoolInitNum  int
	PoolCapacity int
	PoolTTL      time.Duration
	Wrappers     []Wrapper
	CallOptions  CallOptions
	Context      context.Context
	Codecs       map[string]codec.NewCodec

	Broker broker.Broker
}

func NewOptions(options ...Option) Options {
	opts := Options{
		Context:     context.Background(),
		ContentType: DefaultContentType,
		Codecs:      make(map[string]codec.NewCodec),
		CallOptions: CallOptions{
			Retries:        DefaultRetries,
			RequestTimeout: DefaultRequestTimeout,
		},
		PoolInitNum:  5,
		PoolCapacity: 20,
		PoolTTL:      DefaultPoolTTL,
		Broker:       broker.DefaultBroker,
	}

	for _, o := range options {
		o(&opts)
	}

	return opts
}

type PublishOptions struct {
	// Exchange is the routing exchange for the message
	Exchange string
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type MessageOptions struct {
	ContentType string
}

type RequestOptions struct {
	ContentType string
	Stream      bool

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type CallOptions struct {
	// 重试次数
	Retries int
	// 请求超时
	RequestTimeout time.Duration
	Address        string
	//建立连接超时
	DialTimeout time.Duration
	// 中间件
	CallWrappers []CallWrapper
	Context      context.Context
}

func Broker(b broker.Broker) Option {
	return func(o *Options) {
		o.Broker = b
	}
}

// Codec to be used to encode/decode requests for a given content type
func Codec(contentType string, c codec.NewCodec) Option {
	return func(o *Options) {
		o.Codecs[contentType] = c
	}
}

// Default content type of the client
func ContentType(ct string) Option {
	return func(o *Options) {
		o.ContentType = ct
	}
}

// PoolTTL sets the connection pool ttl
func PoolTTL(d time.Duration) Option {
	return func(o *Options) {
		o.PoolTTL = d
	}
}
func Wrap(w Wrapper) Option {
	return func(o *Options) {
		o.Wrappers = append(o.Wrappers, w)
	}
}

// Adds a Wrapper to the list of CallFunc wrappers
func WrapCall(cw ...CallWrapper) Option {
	return func(o *Options) {
		o.CallOptions.CallWrappers = append(o.CallOptions.CallWrappers, cw...)
	}
}

// Number of retries when making the request.
// Should this be a Call Option?
func Retries(i int) Option {
	return func(o *Options) {
		o.CallOptions.Retries = i
	}
}

// Should this be a Call Option?
func RequestTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.CallOptions.RequestTimeout = d
	}
}
func WithExchange(e string) PublishOption {
	return func(o *PublishOptions) {
		o.Exchange = e
	}
}

// PublishContext sets the context in publish options
func PublishContext(ctx context.Context) PublishOption {
	return func(o *PublishOptions) {
		o.Context = ctx
	}
}

// WithAddress sets the remote addresses to use rather than using service discovery
func WithAddress(a string) CallOption {
	return func(o *CallOptions) {
		o.Address = a
	}
}

// WithCallWrapper is a CallOption which adds to the existing CallFunc wrappers
func WithCallWrapper(cw ...CallWrapper) CallOption {
	return func(o *CallOptions) {
		o.CallWrappers = append(o.CallWrappers, cw...)
	}
}

func WithRetries(i int) CallOption {
	return func(o *CallOptions) {
		o.Retries = i
	}
}

func WithMessageContentType(ct string) MessageOption {
	return func(o *MessageOptions) {
		o.ContentType = ct
	}
}

func WithContentType(ct string) RequestOption {
	return func(o *RequestOptions) {
		o.ContentType = ct
	}
}
