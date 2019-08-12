package endpoint

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/payfazz/fazzkit/examples/server/internal/helloworld/model"
)

//Greet dummy greet endpoint for example
func Greet() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*model.Greet)
		if !ok {
			return nil, errors.New("request not recognized")
		}
		return &model.GreetResponse{
			Message: "Hello " + req.Name,
		}, nil
	}
}
