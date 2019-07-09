package server

import (
	"context"
	"encoding/json"
	"net/http"
)

type err interface {
	error() error
}

//EncodeHTTP generate a encode function to encode response to json
func (e *Endpoint) EncodeHTTP() func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		if e, ok := response.(err); ok && e.error() != nil {
			return e.error()
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
		return nil
	}
}

//EncodeGRPC generate a encode function to encode response to proto struct
func (e *Endpoint) EncodeGRPC(model interface{}) func(ctx context.Context, response interface{}) (interface{}, error) {
	return func(ctx context.Context, response interface{}) (interface{}, error) {
		if model == nil {
			return nil, nil
		}

		str, _ := json.Marshal(response)
		_ = json.Unmarshal([]byte(str), model)
		return model, nil
	}
}
