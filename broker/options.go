package broker

import (
	"context"

	"github.com/micro/go-micro/v2/codec"
)

type Options struct {
	Addrs        []string
	Codec        codec.Marshaler
	ErrorHandler Handler
	Context      context.Context
}

type PublishOptions struct {
	Context context.Context
}

type SubscribeOptions struct {
	AutoAck bool
	Queue   string
	Context context.Context
}

type Option func(*Options)

type PublishOption func(*PublishOptions)

func PublishContext(ctx context.Context) PublishOption {
	return func(o *PublishOptions) {
		o.Context = ctx
	}
}

type SubscribeOption func(*SubscribeOptions)

func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	opt := SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Addrs sets the host addresses to be used by the broker
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// Codec sets the codec used for encoding/decoding used where
// a broker does not support headers
func Codec(c codec.Marshaler) Option {
	return func(o *Options) {
		o.Codec = c
	}
}

// DisableAutoAck will disable auto acking of messages
// after they have been handled.
func DisableAutoAck() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = false
	}
}

// ErrorHandler will catch all broker errors that cant be handled
// in normal way, for example Codec errors
func ErrorHandler(h Handler) Option {
	return func(o *Options) {
		o.ErrorHandler = h
	}
}

// Queue sets the name of the queue to share messages on
func Queue(name string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Queue = name
	}
}

// SubscribeContext set context
func SubscribeContext(ctx context.Context) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Context = ctx
	}
}
