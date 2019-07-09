package server

import "github.com/go-kit/kit/endpoint"

//Use add middleware to server
func (e *Endpoint) Use(middleware endpoint.Middleware) *Endpoint {
	e.Middleware = append(e.Middleware, middleware)
	return e
}

//EndpointWithMiddleware wrap kit endpoint with middleware
func (e *Endpoint) EndpointWithMiddleware() endpoint.Endpoint {
	endpointFunc := e.Endpoint()
	for i := 0; i < len(e.Middleware); i++ {
		endpointFunc = e.Middleware[i](endpointFunc)
	}
	return endpointFunc
}
