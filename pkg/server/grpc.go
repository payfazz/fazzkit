package server

import (
	"github.com/go-kit/kit/transport/grpc"
	"github.com/payfazz/fazzkit/pkg/server/logger"
)

//NewGRPCServer create go kit GRPC server
func (e *Endpoint) NewGRPCServer(decodeModel interface{}, encodeModel interface{}, options ...grpc.ServerOption) *grpc.Server {
	options = append(options,
		grpc.ServerErrorLogger(*logger.GetLogger()),
	)
	return grpc.NewServer(e.EndpointWithMiddleware(), e.DecodeGRPC(decodeModel), e.EncodeGRPC(encodeModel), options...)
}
