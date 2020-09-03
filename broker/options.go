package broker

import (
	"context"

	"core/codec"
)

type Options struct {
	Addrs        []string
	Secure       bool
	Codec        codec.Marshaler
	ErrorHandler Handler

	Context context.Context
}

type Option func(*Options)
type SubscribeOption func(*SubscribeOptions)
type PublishOption func(*PublishOptions)

type PublishOptions struct {
	Context context.Context
}

type SubscribeOptions struct {
	AutoAck bool
	GroupId string

	Context context.Context
}
