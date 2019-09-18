package middleware

import (
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/payfazz/fazzkit/server/logger"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
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
