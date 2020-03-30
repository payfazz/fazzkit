package server

import (
	"github.com/go-kit/kit/endpoint"

	"github.com/go-kit/kit/transport/http"
	httpserver "github.com/payfazz/fazzkit/server/http"
	"github.com/payfazz/fazzkit/server/middleware"
)

//HTTPOption server info
type HTTPOption struct {
	DecodeModel interface{}
	Logger      *Logger
	DecodeFunc  httpserver.DecodeFunc
	EncodeFunc  httpserver.EncodeFunc
}

//NewHTTPServer create go kit HTTP server
func NewHTTPServer(e endpoint.Endpoint, httpOpt HTTPOption, options ...http.ServerOption) *http.Server {
	if httpOpt.DecodeFunc == nil {
		httpOpt.DecodeFunc = httpserver.Decode
	}

	if httpOpt.EncodeFunc == nil {
		httpOpt.EncodeFunc = httpserver.Encode
	}

	return newHTTPServer(e, httpOpt, options...)
}

//NewHTTPJSONServer create go kit HTTP server
func NewHTTPJSONServer(e endpoint.Endpoint, httpOpt HTTPOption, options ...http.ServerOption) *http.Server {
	if httpOpt.DecodeFunc == nil {
		httpOpt.DecodeFunc = httpserver.DecodeJSON
	}

	if httpOpt.EncodeFunc == nil {
		httpOpt.EncodeFunc = httpserver.Encode
	}

	return newHTTPServer(e, httpOpt, options...)
}

func newHTTPServer(e endpoint.Endpoint, httpOpt HTTPOption, options ...http.ServerOption) *http.Server {
	middlewares := middleware.Nop()

	if httpOpt.DecodeModel != nil {
		mval := middleware.Validator()
		middlewares = endpoint.Chain(mval)
	}

	if httpOpt.Logger != nil {
		mlog := middleware.LogAndInstrumentation(
			httpOpt.Logger.Logger,
			httpOpt.Logger.Namespace,
			httpOpt.Logger.Subsystem,
			httpOpt.Logger.Action,
			httpOpt.Logger.Domain,
		)
		middlewares = endpoint.Chain(mlog, middlewares)
	}

	return http.NewServer(e, httpOpt.DecodeFunc(httpOpt.DecodeModel), httpOpt.EncodeFunc(), options...)
}
