package grpc

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
)

var error1 = errors.New(`invalid_code`)
var error2 = errors.New(`invalid_code_me`)

func Test_ErrorMapper(t *testing.T) {
	translator := NewErrorMapper()

	translator.RegisterError(error1, codes.FailedPrecondition)
	translator.RegisterError(error2, codes.Unauthenticated)

	grpcError := translator.GetCode(error1)
	if grpcError != codes.FailedPrecondition {
		t.Error(`code_not_match`)
	}

	grpcError = translator.GetCode(error2)
	if grpcError != codes.Unauthenticated {
		t.Error(`code_not_match`)
	}

	grpcError = translator.GetCode(errors.New(`new_error`))
	if grpcError != codes.Internal {
		t.Error(`code_not_match`)
	}
}
