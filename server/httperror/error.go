package httperror

import (
	"context"
	"encoding/json"
	"net/http"

	transportError "github.com/payfazz/fazzkit/server/transport/error"
)

//ErrorWithStatusCode error with http status code
type ErrorWithStatusCode struct {
	Err        error
	StatusCode int
}

func (e *ErrorWithStatusCode) Error() string {
	return e.Err.Error()
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

//EncodeErrorWithInternalCode ...
func EncodeErrorWithInternalCode(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	code := http.StatusInternalServerError
	internalCode := "-1"
	if errWithStatusCode, ok := err.(*ErrorWithStatusCode); ok {
		code = errWithStatusCode.StatusCode
		if errWithInternalCode, ok := errWithStatusCode.Err.(*transportError.ErrorWithInternalCode); ok {
			internalCode = errWithInternalCode.Code
		}
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
		"code":  internalCode,
	})
}
