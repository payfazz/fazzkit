package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/payfazz/fazzkit/server/http"
)

func Get() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var data [][]string

		data = append(data, []string{"1", "2", "3"})
		data = append(data, []string{"4", "5", "6"})

		return http.CSVResponse{
			Filename: "xyz.csv",
			Data:     data,
		}, nil
	}
}
