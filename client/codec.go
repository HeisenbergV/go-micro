package client

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/encoding"
)

type protoCodec struct{}
type wrapCodec struct{ encoding.Codec }

var (
	defaultGRPCCodecs = map[string]encoding.Codec{
		"application/proto":      protoCodec{},
		"application/protobuf":   protoCodec{},
		"application/grpc":       protoCodec{},
		"application/grpc+proto": protoCodec{},
	}
)

func (protoCodec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case proto.Message:
		return proto.Marshal(m)
	}
	return nil, fmt.Errorf("failed to marshal: %v is not type of *bytes.Frame or proto.Message", v)
}

func (protoCodec) Unmarshal(data []byte, v interface{}) error {
	m, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to unmarshal: %v is not type of proto.Message", v)
	}
	return proto.Unmarshal(data, m)
}

func (protoCodec) Name() string {
	return "proto"
}
