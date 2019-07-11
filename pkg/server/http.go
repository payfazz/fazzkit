package server

import (
	"context"
	"encoding/json"

	netHTTP "net/http"

	"github.com/go-kit/kit/transport/http"
	"github.com/payfazz/fazzkit/pkg/server/logger"
)

//NewHTTPServer create go kit HTTP server
func (e *Endpoint) NewHTTPServer(decodeModel interface{}, options ...http.ServerOption) *http.Server {
	options = append(options,
		http.ServerErrorEncoder(encodeError),
		http.ServerErrorLogger(*logger.GetLogger()),
	)
	return http.NewServer(e.EndpointWithMiddleware(), e.DecodeHTTP(decodeModel), e.EncodeHTTP(), options...)
}

func encodeError(_ context.Context, err error, w netHTTP.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	code := netHTTP.StatusInternalServerError
	if sc, ok := err.(*ErrorWithStatusCode); ok {
		code = sc.statusCode
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
