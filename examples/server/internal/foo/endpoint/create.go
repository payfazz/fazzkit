package endpoint

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/payfazz/fazzkit/examples/server/internal/foo/model"
	"github.com/payfazz/fazzkit/server/httperror"
)

//Create dummy create endpoint for example
func Create() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*model.CreateFoo)
		if !ok {
			return nil, &httperror.ErrorWithStatusCode{errors.New("invalid model"), http.StatusInternalServerError}
		}

		fmt.Println("creating object...", input)
		return request, nil
	}
}
