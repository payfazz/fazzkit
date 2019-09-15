package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"

	"github.com/payfazz/fazzkit/server/httperror"
	"github.com/payfazz/fazzkit/server/logger"
	"github.com/payfazz/fazzkit/server/validator"

	"github.com/prometheus/client_golang/prometheus"
)

var loggers = make(map[string]*logger.Logger)
var lock sync.Mutex

//LogAndInstrumentation wrap endpoint function with logger.Log and logger.Instrumentation
//This middleware will measure request_count and request_latency_microseconds
func LogAndInstrumentation(kitLogger log.Logger, namespace string, subsystem string, action string) endpoint.Middleware {
	var logObj logger.Logger

	key := fmt.Sprintf("%s_%s", namespace, subsystem)

	if val, ok := loggers[key]; ok {
		logObj = *val
	} else {
		lock.Lock()
		logObj = logger.New(
			kitprometheus.NewCounterFrom(prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "request_count",
				Help:      "Number of request received.",
			}, []string{"function"}),
			kitprometheus.NewSummaryFrom(prometheus.SummaryOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, []string{"function"}),
			kitLogger,
		)

		loggers[key] = &logObj
		lock.Unlock()
	}

	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return logObj.Instrumentation("function", action, logObj.Log("function", action, f))
	}
}

//Validator wrap endpoint function to execute validator.v9
func Validator() endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			err = validator.DefaultValidator()(request)
			if err != nil {
				return nil, &httperror.ErrorWithStatusCode{
					Err:        err.Error(),
					StatusCode: http.StatusUnprocessableEntity,
				}
			}
			return f(ctx, request)
		}
	}
}

func Nop() endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			return f(ctx, request)
		}
	}
}
