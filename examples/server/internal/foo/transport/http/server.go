package http

import (
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/payfazz/fazzkit/examples/server/internal/foo/endpoint"
	"github.com/payfazz/fazzkit/examples/server/internal/foo/model"
	"github.com/payfazz/fazzkit/server"
)

//MakeHandler make http handler for foo example
func MakeHandler(logger kitlog.Logger, opts ...kithttp.ServerOption) http.Handler {
	e := endpoint.Create()

	serverInfo := server.InfoHTTP{
		DecodeModel: &model.CreateFoo{},
		Logger:      logger,
		Namespace:   "test",
		Subsystem:   "test",
		Action:      "POST",
	}

	return server.NewHTTPServer(e, serverInfo, opts...)
}
