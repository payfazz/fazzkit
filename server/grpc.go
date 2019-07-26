package server

import (
	"github.com/go-kit/kit/endpoint"

	"github.com/go-kit/kit/transport/grpc"
	grpcserver "github.com/payfazz/fazzkit/server/grpc"
)

//NewGRPCServer create go kit GRPC server
func NewGRPCServer(e endpoint.Endpoint, decodeModel interface{}, encodeModel interface{}, options ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(e, grpcserver.Decode(decodeModel), grpcserver.Encode(encodeModel), options...)
}
