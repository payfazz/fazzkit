package middleware

import (
	"context"
	"github.com/payfazz/fazzkit/server/http"

	"github.com/go-kit/kit/endpoint"
)

func HTTPRecoverer(errMapper http.ErrorMapper) endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r, e := f(ctx, request)
			if e == nil {
				return r, e
			}
			httpCode := errMapper.GetCode(e)
			return r, &http.TransportError{
				Err:  e,
				Code: httpCode,
			}
		}
	}
}
