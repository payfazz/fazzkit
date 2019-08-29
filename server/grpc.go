package server

import (
	"github.com/go-kit/kit/endpoint"

	"github.com/go-kit/kit/transport/grpc"
	grpcserver "github.com/payfazz/fazzkit/server/grpc"
	"github.com/payfazz/fazzkit/server/middleware"
)

//GRPCOption server info
type GRPCOption struct {
	DecodeModel interface{}
	EncodeModel interface{}
	Logger      *Logger
}

//NewGRPCServer create go kit GRPC server
func NewGRPCServer(e endpoint.Endpoint, grpcOpt GRPCOption, options ...grpc.ServerOption) *grpc.Server {
	middlewares := middleware.Nop()

	if grpcOpt.DecodeModel != nil {
		mval := middleware.Validator()
		middlewares = endpoint.Chain(mval)
	}

	if grpcOpt.Logger != nil {
		mlog := middleware.LogAndInstrumentation(
			grpcOpt.Logger.Logger,
			grpcOpt.Logger.Namespace,
			grpcOpt.Logger.Subsystem,
			grpcOpt.Logger.Action,
		)
		middlewares = endpoint.Chain(mlog, middlewares)
	}

	e = middlewares(e)

	return grpc.NewServer(
		e,
		grpcserver.Decode(grpcOpt.DecodeModel),
		grpcserver.Encode(grpcOpt.EncodeModel),
		options...,
	)
}
