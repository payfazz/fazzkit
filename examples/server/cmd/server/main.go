package main

import (
	"net/http"
	"os"

	fazzkithttp "github.com/payfazz/fazzkit/server/http"

	"github.com/go-chi/chi"
	"github.com/oklog/oklog/pkg/group"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	foohttp "github.com/payfazz/fazzkit/examples/server/internal/foo/transport/http"
	helloworldgrpc "github.com/payfazz/fazzkit/examples/server/internal/helloworld/transport/grpc"
	helloworldhttp "github.com/payfazz/fazzkit/examples/server/internal/helloworld/transport/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	logger := kitlog.NewLogfmtLogger(os.Stderr)
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
	logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)

	var g group.Group

	g.Add(func() error {
		lis, err := net.Listen("tcp", ":1301")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		helloworldgrpc.Register(s, logger)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		logger.Log("grpc", "1301")
		return nil
	}, func(err error) {
		panic(err)
	})

	g.Add(func() error {
		r := chi.NewRouter()

		r.Handle("/metrics", promhttp.Handler())
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome"))
		})

		httpLogger := kitlog.With(logger, "component", "http")

		r.Mount("/v1", makeHandler(httpLogger))
		return http.ListenAndServe(":1300", r)
	}, func(err error) {
		panic(err)
	})

	_ = logger.Log("run", g.Run())
}

func makeHandler(logger kitlog.Logger) http.Handler {
	r := chi.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(fazzkithttp.EncodeError),
	}

	r.Post("/foo/{id}", foohttp.MakeHandler(logger, opts...).ServeHTTP)
	r.Get("/helloworld/{name}", helloworldhttp.MakeHandler(logger, opts...).ServeHTTP)
	return r
}
