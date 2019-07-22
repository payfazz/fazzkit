package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/iancoleman/strcase"
)

//ParseGRPC parse request proto message to model
func (e *Endpoint) ParseGRPC(ctx context.Context, request interface{}, model interface{}) (interface{}, error) {
	protoMessage := request.(proto.Message)
	m := jsonpb.Marshaler{}
	stringProtoMessage, err := m.MarshalToString(protoMessage)
	if err != nil {
		return nil, err
	}

	mapString := make(map[string]interface{})

	err = json.Unmarshal([]byte(stringProtoMessage), &mapString)
	for key, value := range mapString {
		keySnake := strcase.ToSnake(key)
		mapString[keySnake] = value
	}

	jsonString, err := json.Marshal(mapString)
	err = json.Unmarshal(jsonString, model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

//ParseHTTPJson parse request body (json) to model
func (e *Endpoint) ParseHTTPJson(ctx context.Context, request *http.Request, model interface{}) (interface{}, error) {
	err := json.NewDecoder(request.Body).Decode(model)

	if err != nil {
		return nil, err
	}

	return model, nil
}
