package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
)

//ParseJSON parse request body (json) to model
func ParseJSON(ctx context.Context, request *http.Request, model interface{}) (interface{}, error) {
	err := json.NewDecoder(request.Body).Decode(model)

	if err != nil {
		return nil, err
	}

	return model, nil
}

func ParseURlEncoded(ctx context.Context, request *http.Request, model interface{}) (interface{}, error) {
	err := request.ParseForm()
	if err != nil {
		return nil, err
	}
	decoder := schema.NewDecoder()
	err = decoder.Decode(model, request.PostForm)
	if err != nil {
		return nil, err
	}

	return model, nil
}
