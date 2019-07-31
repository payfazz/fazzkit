package server

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"

	"context"
	"encoding/json"

	netHTTP "net/http"

	"github.com/go-kit/kit/transport/http"
	httpserver "github.com/payfazz/fazzkit/server/http"
	"github.com/payfazz/fazzkit/server/middleware"
	"github.com/payfazz/fazzkit/server/servererror"
)

//InfoHTTP server info
type InfoHTTP struct {
	DecodeModel interface{}
	Logger      log.Logger
	Namespace   string
	Subsystem   string
	Action      string
}

//NewHTTPServer create go kit HTTP server
func NewHTTPServer(e endpoint.Endpoint, info InfoHTTP, options ...http.ServerOption) netHTTP.Handler {
	options = append(options, http.ServerErrorEncoder(encodeError))

	mval := middleware.Validator()
	mlog := middleware.LogAndInstrumentation(info.Logger, info.Namespace, info.Subsystem, info.Action)

	middlewares := endpoint.Chain(mlog, mval)
	e = middlewares(e)

	return http.NewServer(e, httpserver.Decode(info.DecodeModel), httpserver.Encode(), options...)
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
