package endpoint

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/payfazz/fazzkit/examples/server/httperror"
	"github.com/payfazz/fazzkit/examples/server/internal/foo/model"
)

//Create dummy create endpoint for example
func Create() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*model.CreateFoo)
		if !ok {
			return nil, &httperror.ErrorWithStatusCode{"invalid model", http.StatusInternalServerError}
		}

		fmt.Println("creating object...", input)
		return request, nil
	}
}
