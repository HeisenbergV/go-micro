package client

import (
	"context"
	"time"
)

type Client interface {
	Init(...Option) error
	NewRequest(service, endpoint string, req interface{}) Request
	Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error
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

type Response interface {
	Header() map[string]string
	Read() ([]byte, error)
}

// PublishOption used by Publish
type PublishOption func(*PublishOptions)

type Option func(*Options)

type CallOption func(*CallOptions)

var (
	// DefaultClient is a default client to use out of the box
	DefaultClient Client = newRpcClient()
	// DefaultBackoff is the default backoff function for retries
	DefaultBackoff = exponentialBackoff
	// DefaultRetry is the default check-for-retry function for retries
	DefaultRetry = RetryOnError
	// DefaultRetries is the default number of times a request is tried
	DefaultRetries = 1
	// DefaultRequestTimeout is the default request timeout
	DefaultRequestTimeout = time.Second * 5
	// DefaultPoolSize sets the connection pool size
	DefaultPoolSize = 100
	// DefaultPoolTTL sets the connection pool ttl
	DefaultPoolTTL = time.Minute

	// NewClient returns a new client
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

// Creates a new request using the default client. Content Type will
// be set to the default within options and use the appropriate codec
func NewRequest(service, endpoint string, request interface{}, reqOpts ...RequestOption) Request {
	return DefaultClient.NewRequest(service, endpoint, request, reqOpts...)
}

func String() string {
	return DefaultClient.String()
}
