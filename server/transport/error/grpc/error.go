package grpc

import (
	"google.golang.org/grpc/codes"
)

type TransportError struct {
	Error error
	Code  codes.Code
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

func (e *ErrorMapper) RegisterError(err error, grpcCode codes.Code) {
	newTransportError := TransportError{
		Error: err,
		Code:  grpcCode,
	}

	e.Error[err] = &newTransportError
}

func (e *ErrorMapper) GetCode(err error) codes.Code {
	if e.Error[err] == nil {
		return codes.Internal
	}

	return e.Error[err].Code
}
