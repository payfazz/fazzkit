package middleware

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	fazzkithttp "github.com/payfazz/fazzkit/server/http"
	"github.com/payfazz/fazzkit/server/validator"
)

//Validator wrap endpoint function to execute validator.v9
func Validator() endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			err = validator.DefaultValidator()(request)
			if err != nil {
				return nil, &fazzkithttp.TransportError{
					Err:  err,
					Code: http.StatusUnprocessableEntity,
				}
			}
			return f(ctx, request)
		}
	}
}
