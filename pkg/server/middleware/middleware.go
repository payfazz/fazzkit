package middleware

import (
	"github.com/go-kit/kit/endpoint"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/payfazz/fazzkit/pkg/server/logger"

	"github.com/prometheus/client_golang/prometheus"
)

//LogAndInstrumentation wrap function with logger.Log and logger.Instrumentation
func LogAndInstrumentation(namespace string, subsystem string, action string) endpoint.Middleware {
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
		*logger.GetLogger(),
	)

	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return logObj.Instrumentation("method", action, logObj.Log("method", action, f))
	}
}
