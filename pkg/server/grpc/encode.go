package grpc

import (
	"context"
	"encoding/json"
)

//Encode generate a encode function to encode response to proto struct
func Encode(model interface{}) func(ctx context.Context, response interface{}) (interface{}, error) {
	return func(ctx context.Context, response interface{}) (interface{}, error) {
		if model == nil {
			return nil, nil
		}

		str, _ := json.Marshal(response)
		_ = json.Unmarshal([]byte(str), model)
		return model, nil
	}
}
