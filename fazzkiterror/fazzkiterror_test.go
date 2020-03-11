package fazzkiterror_test

import (
	"errors"
	"github.com/payfazz/fazzkit/fazzkiterror"
	"github.com/payfazz/fazzkit/server/grpc"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestGetGRPCStatusCode(t *testing.T) {
	f := fazzkiterror.InternalCodeFactory{
		ServiceCode: 10,
		DomainCode:  2,
	}

	err := &grpc.TransportError{
		Err:  errors.New("test"),
		Code: codes.AlreadyExists,
	}

	err2 := f.New(err)
	statusCode := grpc.GetGRPCStatusCode(err2)
	if statusCode != codes.AlreadyExists {
		t.Fatal("already exists expected")
	}
}
