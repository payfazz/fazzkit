package server

import (
	"github.com/go-kit/kit/transport/grpc"
)

//NewGRPCServer create go kit GRPC server
func (e *Endpoint) NewGRPCServer(decodeModel interface{}, encodeModel interface{}, options ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(e.EndpointWithMiddleware(), e.DecodeGRPC(decodeModel), e.EncodeGRPC(encodeModel), options...)
}
