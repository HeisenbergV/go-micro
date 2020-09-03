package client

import (
	"context"
)

type CallFunc func(ctx context.Context, req Request, rsp interface{}, opts CallOptions) error
type CallWrapper func(CallFunc) CallFunc
type Wrapper func(Client) Client
