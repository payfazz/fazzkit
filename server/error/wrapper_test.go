package error_test

import (
	"errors"
	"fmt"
	servererror "github.com/payfazz/fazzkit/server/error"
	"testing"
)

type Foo struct {
	err error
}

func (f Foo) Error() string {
	return "this is foo: " + f.err.Error()
}

func (f Foo) Wrappee() error {
	return f.err
}

func NewFoo(err error) servererror.Wrapper {
	return &Foo {
		err: err,
	}
}

type Bar struct {
	err error
}

func (f Bar) Error() string {
	return "this is bar: " + f.err.Error()
}

func (f Bar) Wrappee() error {
	return f.err
}

func NewBar(err error) servererror.Wrapper {
	return &Bar {
		err: err,
	}
}


func TestTraversal(t *testing.T) {
	err := NewFoo(NewBar(errors.New("test")))
	recursive(err)
}

func recursive(err error) {
	fmt.Println(err.Error())
	if e, ok := err.(servererror.Wrapper); ok {
		recursive(e.Wrappee())
	}
}