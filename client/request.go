package client

import (
	"fmt"
	"strings"
)

type grpcRequest struct {
	service     string
	method      string
	contentType string
	request     interface{}
}

func methodToGRPC(service, method string) string {
	// no method or already grpc method
	if len(method) == 0 || method[0] == '/' {
		return method
	}

	// assume method is Foo.Bar
	mParts := strings.Split(method, ".")
	if len(mParts) != 2 {
		return method
	}

	if len(service) == 0 {
		return fmt.Sprintf("/%s/%s", mParts[0], mParts[1])
	}

	// return /pkg.Foo/Bar
	return fmt.Sprintf("/%s.%s/%s", service, mParts[0], mParts[1])
}

func newGRPCRequest(service, method string, request interface{}, contentType string) Request {
	return &grpcRequest{
		service:     service,
		method:      method,
		request:     request,
		contentType: contentType,
	}
}

func (g *grpcRequest) Service() string {
	return g.service
}

func (g *grpcRequest) ContentType() string {
	return g.contentType
}

func (g *grpcRequest) Method() string {
	return g.method
}

func (g *grpcRequest) Body() interface{} {
	return g.request
}

func (g *grpcRequest) Endpoint() string {
	return g.method
}
