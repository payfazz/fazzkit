package httperror

import (
	servererror "github.com/payfazz/fazzkit/server/error"
	transporterror "github.com/payfazz/fazzkit/server/transport/error"
	"net/http"
)

var defaultInternalCode = "-1"

func SetDefaultInternalCode(code string) {
	defaultInternalCode = code
}

func getInternalCode(err error) string {
	if e, ok := err.(*transporterror.ErrorWithInternalCode); ok {
		return e.Code
	}

	if e, ok := err.(servererror.Wrapper); ok {
		return getInternalCode(e.Wrappee())
	}

	return defaultInternalCode
}

var defaultStatusCode = http.StatusInternalServerError

func SetDefaultStatusCode(code int) {
	defaultStatusCode = code
}

func getStatusCode(err error) int {
	if e, ok := err.(*ErrorWithStatusCode); ok {
		return e.StatusCode
	}

	if e, ok := err.(servererror.Wrapper); ok {
		return getStatusCode(e.Wrappee())
	}

	return defaultStatusCode
}
