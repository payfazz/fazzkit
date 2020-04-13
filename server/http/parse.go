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

func ParseURlEncoded(ctx context.Context, request *http.Request, model interface{}) (interface{}, error) {
	err := request.ParseForm()
	if err != nil {
		return nil, err
	}
	requestMap := make(map[string]interface{})
	for key, val := range request.Form {
		requestMap[key] = val[0]
	}

	requestByte, err := json.Marshal(requestMap)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(requestByte, model)
	if err != nil {
		return nil, err
	}

	return model, nil
}
