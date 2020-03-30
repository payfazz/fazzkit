package http

import (
	"github.com/payfazz/fazzkit/examples/server/internal/testcsv/endpoint"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/payfazz/fazzkit/server"
	fazzkithttp "github.com/payfazz/fazzkit/server/http"
)

//MakeHandler make http handler for foo example
func MakeHandler(logger kitlog.Logger, opts ...kithttp.ServerOption) http.Handler {
	e := endpoint.Get()

	serverInfo := server.HTTPOption{
		DecodeModel: nil,
		Logger: &server.Logger{
			Logger:    logger,
			Namespace: "test",
			Subsystem: "test",
			Action:    "GET",
		},
		EncodeFunc: fazzkithttp.EncodeCSV,
	}

	return server.NewHTTPServer(e, serverInfo, opts...)
}
