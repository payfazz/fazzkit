package fazzkiterror

import (
	"fmt"
	"runtime"
)

type ErrorWithInternalCode struct {
	ServiceCode int
	DomainCode  int
	Err         error
	Code        func() string
}

func NewErrorWithInternalCode(serviceCode, domainCode int, err error) Wrapper {
	e := &ErrorWithInternalCode{
		ServiceCode: serviceCode,
		DomainCode:  domainCode,
		Err:         err,
	}

	e.Code = e.generateCode()
	return e
}

func (err *ErrorWithInternalCode) Error() string {
	return err.Err.Error()
}

func (err *ErrorWithInternalCode) Wrappee() error {
	return err.Err
}

func (err *ErrorWithInternalCode) generateCode() func() string {
	return func() string {
		_, _, line, _ := runtime.Caller(1)
		return fmt.Sprintf("%02x%02x%04d", err.ServiceCode, err.DomainCode, line)
	}
}
