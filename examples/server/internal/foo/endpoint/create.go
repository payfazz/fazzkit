package endpoint

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

//Create dummy create endpoint for example
func Create() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		fmt.Println("creating object...")
		return request, nil
	}
}
