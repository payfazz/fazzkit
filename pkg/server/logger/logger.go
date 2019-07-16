package logger

import (
	"context"
	"encoding/json"
	"fmt"
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
type Callback (func(request Request) (interface{}, error))

//Logger ...
type Logger struct {
	callUpdate     chan (interface{})
	callError      chan (error)
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	logger         log.Logger
}

//New create gokit layer Logger
func New(counter metrics.Counter, latency metrics.Histogram, logger log.Logger) Logger {
	return Logger{
		callUpdate:     make(chan (interface{})),
		callError:      make(chan (error)),
		requestCount:   counter,
		requestLatency: latency,
		logger:         logger,
	}
}

//Instrumentation ...
func (m Logger) Instrumentation(
	method string,
	action string,
	f func(ctx context.Context, request interface{}) (interface{}, error),
) func(ctx context.Context, request interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		defer func(begin time.Time) {
			m.requestCount.With(method, action).Add(1)
			if err != nil {
				m.requestCount.With(method, fmt.Sprintf("%s_FAILED", action)).Add(1)
				m.requestLatency.With(method, fmt.Sprintf("%s_FAILED", action)).Observe(time.Since(begin).Seconds())
			} else {
				m.requestCount.With(method, fmt.Sprintf("%s_SUCCESS", action)).Add(1)
				m.requestLatency.With(method, fmt.Sprintf("%s_SUCCESS", action)).Observe(time.Since(begin).Seconds())
			}
			m.requestLatency.With(method, action).Observe(time.Since(begin).Seconds())
		}(time.Now())
		return f(ctx, request)
	}
}

//Log ...
func (m Logger) Log(
	method string,
	action string,
	f func(ctx context.Context, request interface{}) (interface{}, error),
) func(ctx context.Context, request interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		defer func(begin time.Time) {
			jsonString, _ := json.Marshal(request)
			m.logger.Log(
				"method", method,
				"action", action,
				"params", jsonString,
				"took", time.Since(begin),
				"err", err,
			)
		}(time.Now())
		return f(ctx, request)
	}
}
