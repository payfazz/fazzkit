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

```
import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/payfazz/fazzkit/server/servererror"
)

type FooModel struct {
    bar int
    baz string
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


```
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
		DecodeModel: &FooModel{},
		Logger:      logger,
		Namespace:   "test",
		Subsystem:   "test",
		Action:      "POST",
	}

	return server.NewHTTPServer(e, serverInfo, opts...)
}
```

<a name="http_transport"/>

### Implement Endpoint to HTTP Transport

Use **NewHTTPServer** method from object created from **NewEndpoint()** function in server package. By default, HTTP decode data based on json tag on models.

| Param                             | Description                                                             |
|-----------------------------------|:------------------------------------------------------------------------|
| decodeModel &lt;interface{}>         | empty decode model, must be an address to struct model tagged with json |
| ...options &lt;ServerOption>         | [go-kit grpc server option](https://godoc.org/github.com/go-kit/kit/transport/grpc#ServerOption) |

#### Example

[https://github.com/payfazz/kitx/blob/master/internal/domain/user/transport/http/server.go](https://github.com/payfazz/kitx/blob/master/internal/domain/user/transport/http/server.go)

```
createEndpoint := endpoint.CreateEndpoint()
createEndpoint.Use(middleware.LogAndInstrumentation(libkitUser, "URL___METHOD", "/user"))

createEndpoint.NewHTTPServer(&model.CreateUser{}, options...)
```

### Decode HTTP data using URL parameter

Use **httpurl** tag on models.

#### Example

```
type User struct {
    ID *string `httpurl:"id" validate:"required" json:"id"`
}
```

<a name="grpc_transport"/>

### Implement Endpoint to GRPC Transport

Use **NewGRPCServer** method from object created from **NewEndpoint()** function in server package. By default, GRPC decode data based on json tag on models.

| Param                             | Description                                                             |
|-----------------------------------|:------------------------------------------------------------------------|
| decodeModel &lt;interface{}>         | empty decode model, must be an address to struct model tagged with json |
| encodeModel &lt;interface{}>         | empty encode model, must be an address to response protobuf struct      |
| ...options &lt;ServerOption>         | [go-kit grpc server option](https://godoc.org/github.com/go-kit/kit/transport/grpc#ServerOption) |

#### Example

[https://github.com/payfazz/kitx/blob/master/internal/domain/user/transport/grpc/server.go](https://github.com/payfazz/kitx/blob/master/internal/domain/user/transport/grpc/server.go)

```
createEndpoint := endpoint.CreateEndpoint()
createEndpoint.Use(middleware.LogAndInstrumentation(libkitUser, "grpc_function", "create"))

createEndpoint.NewGRPCServer(&model.CreateUser{}, &pb.CreateUserResponse{}, options...)
```
<a name="validator"/>

### Validator

By default, endpoint will using [gopkg.in/go-playground/validator.v9](https://gopkg.in/go-playground/validator.v9) on decode struct.


#### Example

[https://github.com/payfazz/kitx/blob/master/internal/domain/user/model/createuser.go](https://github.com/payfazz/kitx/blob/master/internal/domain/user/model/createuser.go)

```
type CreateUser struct {
    Username *string `json:"username" validate:"required"`
    Password *string `json:"password" validate:"required"`
    FooBar   string  `validate:"min=3"`
}
```

<a name="override_validator"/>

### Override Validator

Endpoint validator can be added with another validator functions using **AddValidator(ValidationFunc)**. To reset all existing validators, use **SetValidator(ValidationFunc)**.

```
type ValidationFunc func(req interface{}) error
```

#### Example

```
type Create struct{}

func CreateEndpoint() *server.Endpoint {
    createObj := &Create{}
    return server.NewEndpoint(createObj).SetValidator(validate)
}

func validate(req interface{}) error {
    data := req.(*model.CreateUser)
    if *data.Username == "" {
        return errors.New("username cannot be null")
    }
    if *data.Password == "" {
        return errors.New("password cannot be null")
    }
    return nil
}
```
