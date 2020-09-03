package client

import (
	"context"
	"time"
)

type Client interface {
	Init(...Option) error
	NewRequest(service, method string, req interface{}, reqOpts ...RequestOption) Request
	Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error
	Publish(ctx context.Context, msg Message, opts ...PublishOption) error
	String() string
	Options() Options
}

type Request interface {
	Service() string
	Method() string
	Endpoint() string
	ContentType() string
	Body() interface{}
}

type Message interface {
	Topic() string
	Payload() interface{}
	ContentType() string
}
type Response interface {
	Header() map[string]string
	Read() ([]byte, error)
}

// PublishOption used by Publish
type PublishOption func(*PublishOptions)

type Option func(*Options)

type CallOption func(*CallOptions)

// MessageOption used by NewMessage
type MessageOption func(*MessageOptions)

// RequestOption used by NewRequest
type RequestOption func(*RequestOptions)

var (
	// DefaultClient is a default client to use out of the box
	DefaultClient Client = newRpcClient()
	// DefaultBackoff is the default backoff function for retries
	DefaultRetries        = 1
	DefaultRequestTimeout = time.Second * 5
	DefaultPoolSize       = 100
	DefaultPoolTTL        = time.Minute
	DefaultContentType    = "application/protobuf"

	NewClient func(...Option) Client = newRpcClient
)

// Makes a synchronous call to a service using the default client
func Call(ctx context.Context, request Request, response interface{}, opts ...CallOption) error {
	return DefaultClient.Call(ctx, request, response, opts...)
}

// Publishes a publication using the default client. Using the underlying broker
// set within the options.
func Publish(ctx context.Context, msg Message, opts ...PublishOption) error {
	return DefaultClient.Publish(ctx, msg, opts...)
}

func NewRequest(service, endpoint string, request interface{}, reqOpts ...RequestOption) Request {
	return DefaultClient.NewRequest(service, endpoint, request, reqOpts...)
}

func String() string {
	return DefaultClient.String()
}
