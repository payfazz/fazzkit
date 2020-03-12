package fazzkiterror_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/payfazz/fazzkit/fazzkiterror"
	"github.com/payfazz/fazzkit/server/middleware"
	"os"
	"testing"
)

var f fazzkiterror.InternalCodeFactory

func init() {
	f = fazzkiterror.InternalCodeFactory{
		ServiceCode: 15,
		DomainCode:  3,
	}
}

func TestInternalCode(t *testing.T) {
	err := f.New(errors.New("test"))
	e := err.(*fazzkiterror.ErrorWithInternalCode)

	foo(e)
	fmt.Println(e.Code)
}

func foo(e *fazzkiterror.ErrorWithInternalCode) {
	fmt.Println(e.Code)
}

func TestLog(t *testing.T) {
	logger := log.NewLogfmtLogger(os.Stdout)
	m := middleware.LogAndInstrumentation(logger, "ns", "ss", "ac", "dm")
	_, _ = m(dummyEndpoint)(context.Background(), map[string]interface{}{})
}

func dummyEndpoint(ctx context.Context, request interface{}) (response interface{}, err error) {
	fmt.Println(request)
	err = f.New(errors.New("test_dummy_endpoint"))
	return nil, err
}
