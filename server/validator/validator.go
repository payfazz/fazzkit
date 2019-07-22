package validator

import (
	"sync"

	"github.com/payfazz/go-validator/pkg/validator"
)

//ValidationFunc function used in server validator
type ValidationFunc func(req interface{}) error

var onceVal sync.Once
var val *validator.Validator

//DefaultValidator execute ValidateStruct
func DefaultValidator() ValidationFunc {
	onceVal.Do(func() {
		val = validator.New()
	})

	return func(req interface{}) error {
		return val.ValidateStruct(req)
	}
}
