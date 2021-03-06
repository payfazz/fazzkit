package logger

import (
	"context"
	"encoding/json"
	"github.com/payfazz/fazzkit/fazzkiterror"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	fazzkitgrpc "github.com/payfazz/fazzkit/server/grpc"
	fazzkithttp "github.com/payfazz/fazzkit/server/http"
	"google.golang.org/grpc/codes"
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
				errStatus := "failed"
				if fazzkithttp.HasHTTPTransportError(err) {
					code := fazzkithttp.GetHTTPStatusCode(err)
					if code == http.StatusInternalServerError {
						errStatus = "critical"
					}
				} else {
					code := fazzkitgrpc.GetGRPCStatusCode(err)
					if code == codes.Internal {
						errStatus = "critical"
					}
				}

				labelValues = append(labelValues, "status", errStatus)
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
			kv := make([]interface{}, len(keyvals))
			for i := 0; i < len(keyvals); i++ {
				kv[i] = keyvals[i]
			}

			jsonString, _ := json.Marshal(request)
			kv = append(kv,
				"params", string(jsonString),
				"took", time.Since(begin).String(),
			)

			if nil != err {
				kv = append(kv, "err", err.Error())
				internal := fazzkiterror.GetInternalCode(err)
				if internal != fazzkiterror.DefaultInternalCode {
					kv = append(kv, "err_code", internal)
				}
			}
			_ = m.logger.Log(kv...)
		}(time.Now())
		return f(ctx, request)
	}
}
