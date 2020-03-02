package httperror

import (
	"context"
	"encoding/json"
	"net/http"
)

//ErrorWithStatusCode error with http status code
type ErrorWithStatusCode struct {
	Err        error
	StatusCode int
}

func (e *ErrorWithStatusCode) Error() string {
	return e.Err.Error()
}

func (e *ErrorWithStatusCode) Wrappee() error {
	return e.Err
}

//EncodeError ...
func EncodeError(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	code := http.StatusInternalServerError
	if sc, ok := err.(*ErrorWithStatusCode); ok {
		code = sc.StatusCode
	}

	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

//EncodeErrorWithInternalCode ...
func EncodeErrorWithInternalCode(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	code := getStatusCode(err)
	internalCode := getInternalCode(err)

	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(getErrorMap(err, internalCode))
}

func getErrorMap(err error, internalCode string) map[string]interface{} {
	errString := err.Error()
	var errMap map[string]interface{}
	e := json.Unmarshal([]byte(errString), &errMap)
	if nil != e {
		return map[string]interface{}{
			"error": err.Error(),
			"code":  internalCode,
		}
	}

	errMap["code"] = internalCode

	return errMap
}
