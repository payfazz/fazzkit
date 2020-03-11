package grpc

import (
	"google.golang.org/grpc/codes"
)

type TransportError struct {
	Err  error
	Code codes.Code
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

func (e *ErrorMapper) RegisterError(err error, grpcCode codes.Code) {
	newTransportError := TransportError{
		Err:  err,
		Code: grpcCode,
	}

	e.Error[err] = &newTransportError
}

func (e *ErrorMapper) GetCode(err error) codes.Code {
	if e.Error[err] == nil {
		return codes.Internal
	}

	return e.Error[err].Code
}

//ErrorWithStatusCode error with http status code
type ErrorWithStatusCode struct {
	Err        error
	StatusCode codes.Code
}

func (e *ErrorWithStatusCode) Error() string {
	return e.Err.Error()
}

func (e *ErrorWithStatusCode) Wrappee() error {
	return e.Err
}
