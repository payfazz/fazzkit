package grpc

import (
	"context"
	"encoding/json"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/iancoleman/strcase"
)

//Parse request proto message to model
func Parse(ctx context.Context, request interface{}, model interface{}) (interface{}, error) {
	protoMessage := request.(proto.Message)
	m := jsonpb.Marshaler{}
	stringProtoMessage, err := m.MarshalToString(protoMessage)
	if err != nil {
		return nil, err
	}

	mapString := make(map[string]interface{})

	err = json.Unmarshal([]byte(stringProtoMessage), &mapString)

	if m, ok := model.(proto.Message); ok {
		jsonByte, err := json.Marshal(mapString)
		err = jsonpb.UnmarshalString(string(jsonByte), m)
		if err != nil {
			return nil, err
		}

		return m, nil
	}

	for key, value := range mapString {
		keySnake := strcase.ToSnake(key)
		mapString[keySnake] = value
	}

	jsonByte, err := json.Marshal(mapString)
	err = json.Unmarshal(jsonByte, model)
	if err != nil {
		return nil, err
	}

	return model, nil
}
