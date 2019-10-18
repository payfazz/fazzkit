package logger

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

//Request ...
type Request struct {
	Ctx    context.Context
	Method string
	Action string
	Origin string
	Params interface{}
}

//Callback ...
type Callback func(request Request) (interface{}, error)

//Logger ...
type Logger struct {
	callUpdate     chan interface{}
	callError      chan error
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	logger         log.Logger
}

//New create gokit layer Logger
func New(counter metrics.Counter, latency metrics.Histogram, logger log.Logger) Logger {
	return Logger{
		callUpdate:     make(chan interface{}),
		callError:      make(chan error),
		requestCount:   counter,
		requestLatency: latency,
		logger:         logger,
	}
}

//Instrumentation ...
func (m Logger) Instrumentation(
	f func(ctx context.Context, request interface{}) (interface{}, error),
	keyvals ...interface{},

) func(ctx context.Context, request interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		defer func(begin time.Time) {
			labelValues := make([]string, len(keyvals))
			for i := 0; i < len(keyvals); i++ {
				labelValues[i] = keyvals[i].(string)
			}

			if err != nil {
				labelValues = append(labelValues, "status", "failed")
			} else {
				labelValues = append(labelValues, "status", "success")
			}

			m.requestCount.With(labelValues...).Add(1)
			m.requestLatency.With(labelValues...).Observe(time.Since(begin).Seconds())
		}(time.Now())
		return f(ctx, request)
	}
}

//Log ...
func (m Logger) Log(
	f func(ctx context.Context, request interface{}) (interface{}, error),
	keyvals ...interface{},
) func(ctx context.Context, request interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		defer func(begin time.Time) {
			jsonString, _ := json.Marshal(request)
			keyvals = append(keyvals,
				"params", string(jsonString),
				"took", time.Since(begin).String(),
			)

			if nil != err {
				keyvals = append(keyvals, "err", err.Error())
			}
			_ = m.logger.Log(keyvals...)
		}(time.Now())
		return f(ctx, request)
	}
}
