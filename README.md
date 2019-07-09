# fazzkit

##### Table of Contents  
1. [Server](#server)
    * [Endpoint](#endpoint)
    * [Implement Endpoint to GRPC Transport](#grpc_transport)
    * [Implement Endpoint to HTTP Transport](#http_transport)
    * [Decode Options](#http_decode_options)
    * [Validator](#validator)
    * [Override Validator](#override_validator)

<a name="server"/>

## Server

<a name="endpoint"/>

### Endpoint

Create server endpoint using **NewEndpoint()** function in server package. NewEndpoint parameter must be a struct that implements [endpoint interface](https://github.com/payfazz/kitx/blob/master/pkg/server/server.go). Endpoint interface has abstract method **Endpoint()** that return [go-kit endpoint function](https://godoc.org/github.com/go-kit/kit/endpoint#Endpoint).

#### Example

[https://github.com/payfazz/kitx/blob/master/internal/domain/user/endpoint/create.go](https://github.com/payfazz/kitx/blob/master/internal/domain/user/endpoint/create.go)
```
type Create struct{}

func CreateEndpoint() *server.Endpoint {
    createObj := &Create{}
    return server.NewEndpoint(createObj)
}

func (c *Create) Endpoint() kitEndpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (response interface{}, err error) {
        reqData := request.(*model.CreateUser)

        hashed := *reqData.Password + "_hashed"
        return &model.User{
            Username:  reqData.Username,
            Password:  &hashed,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        }, nil
    }
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

<a name="http_decode_options"/>

### Decode Options

Add option functions to run before decode. Function must implements GRPCDecodeOptions (for GRPC transport) or HTTPDecodeOptions (for HTTP transport).

Available HTTP options:

| Function                          | Description                                                             |
|-----------------------------------|:------------------------------------------------------------------------|
| GetUrlParam(urlParams []string)   | Decode data from URL parameters instead of request JSON body. URL param must be snake_case version from struct attribute name |

#### Example GetUrlParam

Get *foo* attibute on **CreateUser** model from URL parameter.

[https://github.com/payfazz/kitx/blob/master/internal/domain/user/transport/http/server.go](https://github.com/payfazz/kitx/blob/master/internal/domain/user/transport/http/server.go)

```
decodeParam := kitxserver.HTTPDecodeParam{
    Model: &model.CreateUser{},
    Options: []kitxserver.HTTPDecodeOptions{
        kitxserver.GetUrlParam([]string{"foo"}),
    },
}

createEndpoint.NewHTTPServer(decodeParam, opts...)
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
