package middleware

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func Nop() endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			return f(ctx, request)
		}
	}
}
