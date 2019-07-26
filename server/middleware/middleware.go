package middleware

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"

	httpserver "github.com/payfazz/fazzkit/server/http"
	"github.com/payfazz/fazzkit/server/logger"
	"github.com/payfazz/fazzkit/server/validator"

	"github.com/prometheus/client_golang/prometheus"
)

//LogAndInstrumentation wrap endpoint function with logger.Log and logger.Instrumentation
//This middleware will measure request_count and request_latency_microseconds
func LogAndInstrumentation(kitLogger log.Logger, namespace string, subsystem string, action string) endpoint.Middleware {
	logObj := logger.New(
		kitprometheus.NewCounterFrom(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "request_count",
			Help:      "Number of request received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(prometheus.SummaryOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
		kitLogger,
	)

	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return logObj.Instrumentation("method", action, logObj.Log("method", action, f))
	}
}

//Validator wrap endpoint function to execute validator.v9
func Validator() endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			err = validator.DefaultValidator()(request)
			if err != nil {
				return nil, &httpserver.ErrorWithStatusCode{
					Err:        err.Error(),
					StatusCode: http.StatusUnprocessableEntity,
				}
			}
			return f(ctx, request)
		}
	}
}
