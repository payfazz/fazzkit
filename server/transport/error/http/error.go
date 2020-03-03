package http

import (
	transporterror "github.com/payfazz/fazzkit/server/error"
	"net/http"
)

type TransportError struct {
	Error error
	Code  int
}

func (e *TransportError) Wrappee() error {
	return e.Error
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
		Error: err,
		Code:  httpCode,
	}

	e.Error[err] = &newTransportError
}

func (e *ErrorMapper) GetCode(err error) int {
	if e.Error[err] != nil {
		return e.Error[err].Code
	}

	if w, ok := err.(transporterror.Wrapper); ok {
		return e.GetCode(w.Wrappee())
	}

	return http.StatusInternalServerError
}
