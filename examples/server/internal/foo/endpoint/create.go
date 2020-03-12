package endpoint

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/payfazz/fazzkit/examples/server/internal/foo/model"
	fazzkithttp "github.com/payfazz/fazzkit/server/http"
)

//Create dummy create endpoint for example
func Create() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*model.CreateFoo)
		if !ok {
			return nil, &fazzkithttp.TransportError{errors.New("invalid model"), http.StatusInternalServerError}
		}

		fmt.Println("creating object...", input)
		return request, nil
	}
}
