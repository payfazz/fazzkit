package http

import (
	"net/http"

	kitendpoint "github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/payfazz/fazzkit/examples/server/internal/foo/endpoint"
	"github.com/payfazz/fazzkit/examples/server/internal/foo/model"
	"github.com/payfazz/fazzkit/server"
	"github.com/payfazz/fazzkit/server/middleware"
)

//MakeHandler make http handler for foo example
func MakeHandler(logger kitlog.Logger, opts ...kithttp.ServerOption) http.Handler {
	e := endpoint.Create()

	mval := middleware.Validator()
	mlog := middleware.LogAndInstrumentation(logger, "test", "test", "POST")

	middlewares := kitendpoint.Chain(mlog, mval)
	e = middlewares(e)

	return server.NewHTTPServer(e, &model.CreateFoo{}, opts...)
}
