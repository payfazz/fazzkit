package http

import (
	"errors"
	transporterror "github.com/payfazz/fazzkit/server/transport/error"
	"net/http"
	"testing"
)

var error1 = errors.New(`invalid_code`)
var error2 = errors.New(`invalid_code_me`)

func Test_ErrorMapper(t *testing.T) {
	translator := NewErrorMapper()

	translator.RegisterError(error1, http.StatusUnprocessableEntity)
	translator.RegisterError(error2, http.StatusUnauthorized)

	httpError := translator.GetCode(error1)
	if httpError != http.StatusUnprocessableEntity {
		t.Error(`code_not_match`)
	}

	httpError = translator.GetCode(error2)
	if httpError != http.StatusUnauthorized {
		t.Error(`code_not_match`)
	}

	httpError = translator.GetCode(errors.New(`new_error`))
	if httpError != http.StatusInternalServerError {
		t.Error(`code_not_match`)
	}

	err := transporterror.NewRuntimeError(error1, errors.New("new_error"))
	httpError = translator.GetCode(err)
	if httpError != http.StatusUnprocessableEntity {
		t.Error(`code_not_match`)
	}
}
