package server

import (
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/payfazz/fazzkit/pkg/server/logger"
)

var grpconce sync.Once
var grpcLogger *log.Logger

//NewGRPCServer create go kit GRPC server
func (e *Endpoint) NewGRPCServer(decodeModel interface{}, encodeModel interface{}, options ...grpc.ServerOption) *grpc.Server {
	grpconce.Do(func() {
		logObj := logger.GetLogger()
		_grpcLogger := log.With(*logObj, "component", "grpc")
		grpcLogger = &_grpcLogger
	})

	return grpc.NewServer(e.EndpointWithMiddleware(), e.DecodeGRPC(decodeModel), e.EncodeGRPC(encodeModel), options...)
}
