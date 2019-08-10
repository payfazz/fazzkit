package server

import (
	"github.com/go-kit/kit/endpoint"

	"context"
	"encoding/json"

	netHTTP "net/http"

	"github.com/go-kit/kit/transport/http"
	httpserver "github.com/payfazz/fazzkit/server/http"
	"github.com/payfazz/fazzkit/server/middleware"
	"github.com/payfazz/fazzkit/server/servererror"
)

//HTTPOption server info
type HTTPOption struct {
	DecodeModel interface{}
	Logger      *Logger
}

//NewHTTPServer create go kit HTTP server
func NewHTTPServer(e endpoint.Endpoint, httpOpt HTTPOption, options ...http.ServerOption) netHTTP.Handler {
	options = append(options, http.ServerErrorEncoder(encodeError))

	mval := middleware.Validator()
	middlewares := endpoint.Chain(mval)

	if httpOpt.Logger != nil {
		mlog := middleware.LogAndInstrumentation(
			httpOpt.Logger.Logger,
			httpOpt.Logger.Namespace,
			httpOpt.Logger.Subsystem,
			httpOpt.Logger.Action,
		)
		middlewares = endpoint.Chain(mlog, middlewares)
	}

	e = middlewares(e)

	return http.NewServer(e, httpserver.Decode(httpOpt.DecodeModel), httpserver.Encode(), options...)
}

func encodeError(_ context.Context, err error, w netHTTP.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	code := netHTTP.StatusInternalServerError
	if sc, ok := err.(*servererror.ErrorWithStatusCode); ok {
		code = sc.StatusCode
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
