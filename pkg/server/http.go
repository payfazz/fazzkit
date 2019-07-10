package server

import (
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/http"
	"github.com/payfazz/fazzkit/pkg/server/logger"
)

var httponce sync.Once
var httpLogger *log.Logger

//NewHTTPServer create go kit HTTP server
func (e *Endpoint) NewHTTPServer(decodeModel interface{}, options ...http.ServerOption) *http.Server {
	httponce.Do(func() {
		logObj := logger.GetLogger()
		_httpLogger := log.With(*logObj, "component", "http")
		httpLogger = &_httpLogger
	})

	return http.NewServer(e.EndpointWithMiddleware(), e.DecodeHTTP(decodeModel), e.EncodeHTTP(), options...)
}
