package server

import (
	"github.com/go-kit/kit/endpoint"
)

//Interface describe endpoint methods.
//Endpoint return go-kit endpoint.
type Interface interface {
	Endpoint() endpoint.Endpoint
}

//ValidationFunc function used in server validator
type ValidationFunc func(req interface{}) error

//Endpoint struct must use Interface. Model used for transport decode method
type Endpoint struct {
	Interface

	Validators []ValidationFunc
	Middleware []endpoint.Middleware
}

//NewEndpoint create new endpoint
//validation will be injected to transport decode method
func NewEndpoint(e Interface) *Endpoint {
	endpoint := &Endpoint{e, []ValidationFunc{defaultValidator()}, []endpoint.Middleware{}}
	return endpoint
}
