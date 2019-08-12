//go:generate protoc -I ../../../../pkg/proto/helloworld ../../../../pkg/proto/helloworld/helloworld.proto --go_out=plugins=grpc:../../../../pkg/proto/helloworld
package grpc

import (
	"context"

	kitlog "github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	"github.com/payfazz/fazzkit/examples/server/internal/helloworld/endpoint"
	"github.com/payfazz/fazzkit/examples/server/internal/helloworld/model"
	pb "github.com/payfazz/fazzkit/examples/server/pkg/proto/helloworld"

	"github.com/payfazz/fazzkit/server"
)

//grpcServer ...
type grpcServer struct {
	sayHello kitgrpc.Handler
}

//SayHello ...
func (s *grpcServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	_, rep, err := s.sayHello.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.HelloReply), nil
}

//Register ...
func Register(s *grpc.Server, logger kitlog.Logger, opts ...kitgrpc.ServerOption) {
	sayHello := endpoint.Greet()

	serverInfo := server.GRPCOption{
		DecodeModel: &model.Greet{},
		EncodeModel: &pb.HelloReply{},
		Logger: &server.Logger{
			Logger:    logger,
			Namespace: "test",
			Subsystem: "test",
			Action:    "POST",
		},
	}

	pb.RegisterGreeterServer(s, &grpcServer{
		sayHello: server.NewGRPCServer(sayHello, serverInfo),
	})
}
