package grpc

import (
	"github.com/payfazz/fazzkit/fazzkiterror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

var defaultGRPCStatusCode = codes.Unknown

func SetDefaultGRPCStatusCode(code codes.Code) {
	defaultGRPCStatusCode = code
}

func GetGRPCStatusCode(err error) codes.Code {
	if e, ok := err.(*TransportError); ok {
		return e.Code
	}

	if e, ok := err.(fazzkiterror.Wrapper); ok {
		return GetGRPCStatusCode(e.Wrappee())
	}

	if se, ok := err.(interface {
		GRPCStatus() *status.Status
	}); ok {
		return se.GRPCStatus().Code()
	}

	return defaultGRPCStatusCode
}
