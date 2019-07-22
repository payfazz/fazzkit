package http

import (
	"context"
	"encoding/json"
	"net/http"
)

//ParseJSON parse request body (json) to model
func ParseJSON(ctx context.Context, request *http.Request, model interface{}) (interface{}, error) {
	err := json.NewDecoder(request.Body).Decode(model)

	if err != nil {
		return nil, err
	}

	return model, nil
}
