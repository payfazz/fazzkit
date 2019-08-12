# fazzkit

##### Table of Contents  
1. [Server](#server)
    * [Endpoint](#endpoint)
    * [Implement Endpoint to HTTP Transport](#http_transport)
    * [Implement Endpoint to GRPC Transport](#grpc_transport)
    * [Validator](#validator)
    * [Override Validator](#override_validator)

<a name="server"/>

## Server

<a name="endpoint"/>

### Endpoint

Define a [go-kit endpoint function](https://godoc.org/github.com/go-kit/kit/endpoint#Endpoint).
Define request model tagged with json.

```
import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/payfazz/fazzkit/server/servererror"
)

type FooModel struct {
	bar int    `json:bar`
	baz string `json:baz`
}

func Foo() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*FooModel)
		if !ok {
			return nil, &servererror.ErrorWithStatusCode{"invalid model", http.StatusInternalServerError}
		}

		fmt.Println("processing object...", input)
		return request, nil
	}
}
```

<a name="http_transport"/>

### Implement Endpoint to HTTP Transport

Use **NewHTTPServer** function from server package.

| Param                             | Description                                                             |
|-----------------------------------|:------------------------------------------------------------------------|
| Endpoint                          | [go-kit Endpoint](https://godoc.org/github.com/go-kit/kit/endpoint#Endpoint) |
| HTTPOption                        | fazzkit HTTPOption |
| ...ServerOption                   | [go-kit grpc server option](https://godoc.org/github.com/go-kit/kit/transport/http#ServerOption) |

Put your decode model in fazzkit HTTPOption. This model will be used for decoding and validating request from HTTP. By default, your data decoded from json body.

```
import (
	"net/http"

	"github.com/payfazz/fazzkit/server"
)

//MakeHandler make http handler for foo example
func MakeHandler() http.Handler {
	e := Foo()

	httpOpt := server.HTTPOption{
		DecodeModel: &model.FooModel{},
	}

	return server.NewHTTPServer(e, httpOpt)
}
```

Add fazzkit logger to HTTPOption to measure request count and request latency with prometheus.

```
import (
	"net/http"

	"github.com/payfazz/fazzkit/server"
	"github.com/go-kit/kit/log"
)

//MakeHandler make http handler for foo example
func MakeHandler(logger log.Logger) http.Handler {
	e := Foo()

	serverInfo := server.HTTPOption{
		DecodeModel: &model.CreateFoo{},
		Logger: &server.Logger{
			Logger:    logger,
			Namespace: "test",
			Subsystem: "foo",
			Action:    "GET",
		},
	}

	return server.NewHTTPServer(e, httpOpt)
}
```

Add some go-kit server options when needed.

```
import (
	"net/http"

	"github.com/payfazz/fazzkit/server"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

//MakeHandler make http handler for foo example
func MakeHandler(logger log.Logger) http.Handler {
	e := Foo()

	serverInfo := server.HTTPOption{
		DecodeModel: &model.CreateFoo{},
		Logger: &server.Logger{
			Logger:    logger,
			Namespace: "test",
			Subsystem: "foo",
			Action:    "GET",
		},
	}

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
	}

	return server.NewHTTPServer(e, httpOpt, opts)
}
```

### Decode HTTP data using URL parameter

Use **httpurl** tag on models.

```
type FooModel struct {
	bar int    `json:bar`
	baz string `json:baz`
	ID  string `httpurl:id`
}
```

<a name="grpc_transport"/>

### Implement Endpoint to GRPC Transport

Use **NewGRPCServer** function from server package.

| Param                             | Description                                                             |
|-----------------------------------|:------------------------------------------------------------------------|
| Endpoint                          | [go-kit Endpoint](https://godoc.org/github.com/go-kit/kit/endpoint#Endpoint) |
| GRPCOption                        | fazzkit GRPCOption |
| ...ServerOption                   | [go-kit grpc server option](https://godoc.org/github.com/go-kit/kit/transport/http#ServerOption) |

[Example](https://github.com/payfazz/fazzkit/blob/master/examples/server/internal/helloworld/transport/grpc/server.go)

<a name="validator"/>

### Validator

By default, endpoint will using [gopkg.in/go-playground/validator.v9](https://gopkg.in/go-playground/validator.v9) on decode struct.

```
type CreateUser struct {
    Username *string `json:"username" validate:"required"`
    Password *string `json:"password" validate:"required"`
    FooBar   string  `validate:"min=3"`
}
```
