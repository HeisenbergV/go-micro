package micro

import (
	"context"

	"core/broker"
	"core/client"
	"core/server"
)

type Options struct {
	Broker broker.Broker
	Client client.Client
	Server server.Server

	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error

	Context context.Context

	Signal bool
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Server:  server.DefaultServer,
		Context: context.Background(),
		Signal:  true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Broker to be used for service
func Broker(b broker.Broker) Option {
	return func(o *Options) {
		o.Broker = b
		// Update Client and Server
		o.Client.Init(client.Broker(b))
		o.Server.Init(server.Broker(b))
	}
}

// Client to be used for service
func Client(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service and for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// HandleSignal toggles automatic installation of the signal handler that
// traps TERM, INT, and QUIT.  Users of this feature to disable the signal
// handler, should control liveness of the service through the context.
func HandleSignal(b bool) Option {
	return func(o *Options) {
		o.Signal = b
	}
}

// Server to be used for service
func Server(s server.Server) Option {
	return func(o *Options) {
		o.Server = s
	}
}

// Address sets the address of the server
func Address(addr string) Option {
	return func(o *Options) {
		o.Server.Init(server.Address(addr))
	}
}

// Name of the service
func Name(n string) Option {
	return func(o *Options) {
		o.Server.Init(server.Name(n))
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.Server.Init(server.Version(v))
	}
}

// Metadata associated with the service
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.Server.Init(server.Metadata(md))
	}
}
func WrapClient(w ...client.Wrapper) Option {
	return func(o *Options) {
		// apply in reverse
		for i := len(w); i > 0; i-- {
			o.Client = w[i-1](o.Client)
		}
	}
}

// WrapCall is a convenience method for wrapping a Client CallFunc
func WrapCall(w ...client.CallWrapper) Option {
	return func(o *Options) {
		o.Client.Init(client.WrapCall(w...))
	}
}

// WrapHandler adds a handler Wrapper to a list of options passed into the server
func WrapHandler(w ...server.HandlerWrapper) Option {
	return func(o *Options) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapHandler(wrap))
		}

		// Init once
		o.Server.Init(wrappers...)
	}
}

// WrapSubscriber adds a subscriber Wrapper to a list of options passed into the server
func WrapSubscriber(w ...server.SubscriberWrapper) Option {
	return func(o *Options) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapSubscriber(wrap))
		}

		// Init once
		o.Server.Init(wrappers...)
	}
}

// Before and Afters

// BeforeStart run funcs before service starts
func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

// BeforeStop run funcs before service stops
func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

// AfterStart run funcs after service starts
func AfterStart(fn func() error) Option {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

// AfterStop run funcs after service stops
func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}
