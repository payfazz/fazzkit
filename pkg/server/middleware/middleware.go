package middleware

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/payfazz/kitx/pkg/thunk/logger"
)

func LogAndInstrumentation(logger logger.ILogger, method string, action string) endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return logger.LogAndInstrumentation(method, action, f)
	}
}
