package http

import "encoding/json"

type JSONError struct {
	data map[string]interface{}
}

func (err *JSONError) Error() string {
	data, _ := json.Marshal(err.data)
	return string(data)
}
