package http

import (
	"context"
	"encoding/json"
	"net/http"
)

type err interface {
	error() error
}

type EncodeFunc func() func(ctx context.Context, w http.ResponseWriter, response interface{}) error

//Encode generate a encode function to encode response to json
func Encode() func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		if e, ok := response.(err); ok && e.error() != nil {
			return e.error()
		}

		// set status code to 204 when the response is nil
		if response == nil {
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode("")
			return nil
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
		return nil
	}
}
