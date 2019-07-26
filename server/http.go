package server

import (
	"github.com/go-kit/kit/endpoint"

	"context"
	"encoding/json"

	netHTTP "net/http"

	"github.com/go-kit/kit/transport/http"
	httpserver "github.com/payfazz/fazzkit/server/http"
)

//NewHTTPServer create go kit HTTP server
func NewHTTPServer(e endpoint.Endpoint, decodeModel interface{}, options ...http.ServerOption) netHTTP.HandlerFunc {
	options = append(options, http.ServerErrorEncoder(encodeError))
	return http.NewServer(e, httpserver.Decode(decodeModel), httpserver.Encode(), options...).ServeHTTP
}

func encodeError(_ context.Context, err error, w netHTTP.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	code := netHTTP.StatusInternalServerError
	if sc, ok := err.(*httpserver.ErrorWithStatusCode); ok {
		code = sc.StatusCode
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
