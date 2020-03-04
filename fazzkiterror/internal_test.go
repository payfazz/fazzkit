package fazzkiterror_test

import (
	"errors"
	"fmt"
	"github.com/payfazz/fazzkit/fazzkiterror"
	"testing"
)

func TestInternalCode(t *testing.T) {
	f := fazzkiterror.InternalCodeFactory{
		ServiceCode: 15,
		DomainCode:  3,
	}

	err := f.New(errors.New("test"))
	e := err.(*fazzkiterror.ErrorWithInternalCode)

	foo(e)
	fmt.Println(e.Code)
}

func foo(e *fazzkiterror.ErrorWithInternalCode) {
	fmt.Println(e.Code)
}
