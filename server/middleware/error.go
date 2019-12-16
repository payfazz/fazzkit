package middleware

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/payfazz/fazzkit/server/httperror"
	errorgrpc "github.com/payfazz/fazzkit/server/transport/error/grpc"
	errorhttp "github.com/payfazz/fazzkit/server/transport/error/http"
	"google.golang.org/grpc/status"
)

func TranslateErrorHTTP(errMapper errorhttp.ErrorMapper) endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r, e := f(ctx, request)
			if e == nil {
				return r, e
			}
			httpCode := errMapper.GetCode(e)
			return r, &httperror.ErrorWithStatusCode{
				Err:        e.Error(),
				StatusCode: httpCode,
			}
		}
	}
}

func TranslateErrorGRPC(errMapper errorgrpc.ErrorMapper) endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r, e := f(ctx, request)
			if e == nil {
				return r, e
			}
			grpcCode := errMapper.GetCode(e)
			return r, status.Error(grpcCode, e.Error())
		}
	}
}
