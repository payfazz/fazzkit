package middleware

import (
	"github.com/go-kit/kit/endpoint"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"

	"github.com/go-kit/kit/log"
	"github.com/payfazz/fazzkit/server/logger"
	"github.com/prometheus/client_golang/prometheus"
)

//LogAndInstrumentation wrap function with logger.Log and logger.Instrumentation
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
