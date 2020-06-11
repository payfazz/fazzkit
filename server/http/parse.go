package http

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
)

//ParseJSON parse request body (json) to model
func ParseJSON(ctx context.Context, request *http.Request, model interface{}) (interface{}, error) {
	err := json.NewDecoder(request.Body).Decode(model)

	if err != nil {
		return nil, err
	}

	return model, nil
}

//ParseCSV ...
func ParseCSV(ctx context.Context, request *http.Request, model interface{}) (interface{}, error) {
	typ := reflect.TypeOf(model).Elem()
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("csv")
		if tag == "" {
			continue
		}

		data, err := parseCSV(ctx, request, tag)
		if nil != err {
			return model, err
		}

		reflect.ValueOf(model).Elem().Field(i).Set(reflect.ValueOf(data))
	}
	return model, nil
}

func parseCSV(ctx context.Context, request *http.Request, key string) ([]map[string]interface{}, error) {
	f, _, err := request.FormFile(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make([]map[string]interface{}, 0)

	reader := csv.NewReader(f)
	keys, err := reader.Read()
	if nil != err {
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		for i, r := range record {
			row := make(map[string]interface{})
			row[keys[i]] = r
			result = append(result, row)
		}
	}

	return result, nil
}

func ParseURlEncoded(ctx context.Context, request *http.Request, model interface{}) (interface{}, error) {
	err := request.ParseForm()
	if err != nil {
		return nil, err
	}
	requestMap := make(map[string]interface{})
	for key, val := range request.Form {
		if len(val) > 1 {
			requestMap[key] = val
			continue
		}
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
