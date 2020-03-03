package middleware

import (
	"context"
	"github.com/payfazz/fazzkit/server/grpc"
	"github.com/payfazz/fazzkit/server/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/payfazz/fazzkit/server/httperror"
	"google.golang.org/grpc/status"
)

func TranslateErrorHTTP(errMapper http.ErrorMapper) endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r, e := f(ctx, request)
			if e == nil {
				return r, e
			}
			httpCode := errMapper.GetCode(e)
			return r, &httperror.ErrorWithStatusCode{
				Err:        e,
				StatusCode: httpCode,
			}
		}
	}
}

func TranslateErrorGRPC(errMapper grpc.ErrorMapper) endpoint.Middleware {
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
