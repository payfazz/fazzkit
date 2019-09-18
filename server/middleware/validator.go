package middleware

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/payfazz/fazzkit/server/httperror"
	"github.com/payfazz/fazzkit/server/validator"
	"net/http"
)

//Validator wrap endpoint function to execute validator.v9
func Validator() endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			err = validator.DefaultValidator()(request)
			if err != nil {
				return nil, &httperror.ErrorWithStatusCode{
					Err:        err.Error(),
					StatusCode: http.StatusUnprocessableEntity,
				}
			}
			return f(ctx, request)
		}
	}
}
