package fazzkiterror

import (
	fazzkithttp "github.com/payfazz/fazzkit/server/http"
	fazzkitgrpc "github.com/payfazz/fazzkit/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

var defaultInternalCode = "-1"

func SetDefaultInternalCode(code string) {
	defaultInternalCode = code
}

func GetInternalCode(err error) string {
	if e, ok := err.(*ErrorWithInternalCode); ok {
		return e.Code
	}

	if e, ok := err.(Wrapper); ok {
		return GetInternalCode(e.Wrappee())
	}

	return defaultInternalCode
}

func GetDomainError(err error) string {
	if e, ok := err.(Wrapper); ok {
		return GetDomainError(e.Wrappee())
	}

	return err.Error()
}

var defaultHTTPStatusCode = http.StatusInternalServerError

func SetDefaultHTTPStatusCode(code int) {
	defaultHTTPStatusCode = code
}

func GetHTTPStatusCode(err error) int {
	if e, ok := err.(*fazzkithttp.ErrorWithStatusCode); ok {
		return e.StatusCode
	}

	if e, ok := err.(Wrapper); ok {
		return GetHTTPStatusCode(e.Wrappee())
	}

	return defaultHTTPStatusCode
}

var defaultGRPCStatusCode = codes.Unknown

func SetDefaultGRPCStatusCode(code codes.Code) {
	defaultGRPCStatusCode = code
}

func GetGRPCStatusCode(err error) codes.Code {
	if e, ok := err.(*fazzkitgrpc.ErrorWithStatusCode); ok {
		return e.StatusCode
	}

	if e, ok := err.(Wrapper); ok {
		return GetGRPCStatusCode(e.Wrappee())
	}

	if se, ok := err.(interface {
		GRPCStatus() *status.Status
	}); ok {
		return se.GRPCStatus().Code()
	}

	return defaultGRPCStatusCode
}