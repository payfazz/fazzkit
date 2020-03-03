package fazzkiterror_test

import (
	"errors"
	"fmt"
	"github.com/payfazz/fazzkit/fazzkiterror"
	"testing"
)

func TestInternalCode(t *testing.T) {
	err := fazzkiterror.NewErrorWithInternalCode(15, 3, errors.New("test"))

	e := err.(*fazzkiterror.ErrorWithInternalCode)

	foo(e)
	bar(e.Code)
	fmt.Println(e.Code())
}

func foo(e *fazzkiterror.ErrorWithInternalCode) {
	fmt.Println(e.Code())
}

func bar(f func() string) {
	fmt.Println(f())
	baz(f)
}

func baz(f func() string) {
	fmt.Println(f())
}