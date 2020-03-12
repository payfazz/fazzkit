package http

import (
	"context"
	"encoding/json"
	"github.com/payfazz/fazzkit/fazzkiterror"
	"net/http"
)

type TransportError struct {
	Err  error
	Code int
}

func (e *TransportError) Error() string {
	return e.Err.Error()
}

func (e *TransportError) Wrappee() error {
	return e.Err
}

type ErrorMapper struct {
	Error map[error]*TransportError
}

func NewErrorMapper() *ErrorMapper {
	return &ErrorMapper{
		Error: make(map[error]*TransportError),
	}
}

func (e *ErrorMapper) RegisterError(err error, httpCode int) {
	newTransportError := TransportError{
		Err:  err,
		Code: httpCode,
	}

	e.Error[err] = &newTransportError
}

func (e *ErrorMapper) GetCode(err error) int {
	if e.Error[err] != nil {
		return e.Error[err].Code
	}

	if w, ok := err.(fazzkiterror.Wrapper); ok {
		return e.GetCode(w.Wrappee())
	}

	return http.StatusInternalServerError
}

//EncodeError ...
func EncodeError(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	code := http.StatusInternalServerError
	if sc, ok := err.(*TransportError); ok {
		code = sc.Code
	}

	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

//EncodeErrorWithInternalCode ...
func EncodeErrorWithInternalCode(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	code := GetHTTPStatusCode(err)
	internalCode := fazzkiterror.GetInternalCode(err)

	w.WriteHeader(code)

	errorString := fazzkiterror.GetDomainError(err)
	errMap := getErrorMap(err)
	errMap["code"] = internalCode
	if _, ok := errMap["error"]; !ok {
		errMap["error"] = errorString
	}

	_ = json.NewEncoder(w).Encode(errMap)
}

func getErrorMap(err error) map[string]interface{} {
	errString := err.Error()
	var errMap map[string]interface{}
	e := json.Unmarshal([]byte(errString), &errMap)
	if nil == e {
		return errMap
	}

	return map[string]interface{}{
		"error": errString,
	}
}


var defaultHTTPStatusCode = http.StatusInternalServerError

func SetDefaultHTTPStatusCode(code int) {
	defaultHTTPStatusCode = code
}

func GetHTTPStatusCode(err error) int {
	if e, ok := err.(*TransportError); ok {
		return e.Code
	}

	if e, ok := err.(fazzkiterror.Wrapper); ok {
		return GetHTTPStatusCode(e.Wrappee())
	}

	return defaultHTTPStatusCode
}

func HasHTTPTransportError(err error) bool {
	if _, ok := err.(*TransportError); ok {
		return true
	}

	if e, ok := err.(fazzkiterror.Wrapper); ok {
		return HasHTTPTransportError(e.Wrappee())
	}

	return false
}
