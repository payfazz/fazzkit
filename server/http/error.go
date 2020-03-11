package http

import (
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
