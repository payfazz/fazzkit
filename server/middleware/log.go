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
func LogAndInstrumentation(kitLogger log.Logger, namespace, subsystem, action, domain string) endpoint.Middleware {
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
			}, []string{"function", "status", "domain"}),
			kitprometheus.NewHistogramFrom(prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "request_latency",
				Help:      "Total duration of requests in seconds.",
			}, []string{"function", "status", "domain"}),
			kitLogger,
		)

		loggers[key] = &logObj
		lock.Unlock()
	}

	return func(f endpoint.Endpoint) endpoint.Endpoint {
		keyvals := make([]interface{}, 0)
		keyvals = append(keyvals,
			"function", action,
			"domain", domain,
		)
		return logObj.Instrumentation(logObj.Log(f, keyvals...), keyvals...)
	}
}
