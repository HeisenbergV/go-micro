// Package micro is a pluggable framework for microservices
package micro

import (
	"context"

	"core/client"
	"core/server"
)

type serviceKey struct{}

type Service interface {
	Name() string
	Init(...Option)
	Options() Options
	Client() client.Client
	Server() server.Server
	Run() error
	String() string
}

type Event interface {
	Publish(ctx context.Context, msg interface{}, opts ...client.PublishOption) error
}

// Type alias to satisfy the deprecation
type Publisher = Event

type Option func(*Options)

var (
	HeaderPrefix = "TSP-"
)

func NewService(opts ...Option) Service {
	return newService(opts...)
}

func FromContext(ctx context.Context) (Service, bool) {
	s, ok := ctx.Value(serviceKey{}).(Service)
	return s, ok
}

func NewContext(ctx context.Context, s Service) context.Context {
	return context.WithValue(ctx, serviceKey{}, s)
}

func RegisterHandler(s server.Server, h interface{}, opts ...server.HandlerOption) error {
	return s.Handle(s.NewHandler(h, opts...))
}
