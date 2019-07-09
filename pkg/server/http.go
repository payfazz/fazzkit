package server

import (
	"github.com/go-kit/kit/transport/http"
)

//NewHTTPServer create go kit HTTP server
func (e *Endpoint) NewHTTPServer(decodeModel interface{}, options ...http.ServerOption) *http.Server {
	return http.NewServer(e.EndpointWithMiddleware(), e.DecodeHTTP(decodeModel), e.EncodeHTTP(), options...)
}
