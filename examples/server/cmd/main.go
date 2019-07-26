package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	foohttp "github.com/payfazz/fazzkit/examples/server/internal/foo/transport/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	r := chi.NewRouter()

	r.Handle("/metrics", promhttp.Handler())
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	logger := kitlog.NewLogfmtLogger(os.Stderr)
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
	logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)

	httpLogger := kitlog.With(logger, "component", "http")

	r.Mount("/v1", makeHandler(httpLogger))

	http.ListenAndServe(":1300", r)
}

func makeHandler(logger kitlog.Logger) http.Handler {
	r := chi.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
	}

	r.Post("/foo/{id}", foohttp.MakeHandler(logger, opts...))
	return r
}
