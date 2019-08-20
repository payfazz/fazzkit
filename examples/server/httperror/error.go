package httperror

import (
	"context"
	"encoding/json"
	"net/http"
)

//ErrorWithStatusCode error with http status code
type ErrorWithStatusCode struct {
	Err        string
	StatusCode int
}

func (e *ErrorWithStatusCode) Error() string {
	return e.Err
}

//EncodeError ...
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	code := http.StatusInternalServerError
	if sc, ok := err.(*ErrorWithStatusCode); ok {
		code = sc.StatusCode
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
