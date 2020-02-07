package http

import (
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
	if e.Error[err] == nil {
		return http.StatusInternalServerError
	}

	return e.Error[err].Code
}
