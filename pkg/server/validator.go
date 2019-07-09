package server

import (
	"sync"

	"github.com/payfazz/go-validator/pkg/validator"
)

var onceVal sync.Once
var val *validator.Validator

//SetValidator override existing validators
func (e *Endpoint) SetValidator(validator ValidationFunc) *Endpoint {
	e.Validators = []ValidationFunc{validator}
	return e
}

//AddValidator add validator func to execute in linear order
func (e *Endpoint) AddValidator(validator ValidationFunc) *Endpoint {
	e.Validators = append(e.Validators, validator)
	return e
}

//Validate execute all validator func in linear order
func (e *Endpoint) Validate(req interface{}) error {
	var err error
	for i := 0; i < len(e.Validators); i++ {
		err = e.Validators[i](req)
		if err != nil {
			return err
		}
	}
	return nil
}

func defaultValidator() ValidationFunc {
	onceVal.Do(func() {
		val = validator.New()
	})

	return func(req interface{}) error {
		return val.ValidateStruct(req)
	}
}
