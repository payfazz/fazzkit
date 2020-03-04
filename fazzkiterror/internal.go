package fazzkiterror

import (
	"fmt"
	"runtime"
)

type ErrorWithInternalCode struct {
	ServiceCode int
	DomainCode  int
	Err         error
	Code        string
}

type InternalCodeFactory struct {
	ServiceCode int
	DomainCode  int
}

func (factory *InternalCodeFactory) New(err error) Wrapper {
	w := newErrorWithInternalCode(factory.ServiceCode, factory.DomainCode, err)
	e := w.(*ErrorWithInternalCode)

	_, _, line, _ := runtime.Caller(1)
	e.Code = fmt.Sprintf("%02x%02x%04d", factory.ServiceCode, factory.DomainCode, line)

	return e
}

func newErrorWithInternalCode(serviceCode, domainCode int, err error) Wrapper {
	e := &ErrorWithInternalCode{
		ServiceCode: serviceCode,
		DomainCode:  domainCode,
		Err:         err,
	}
	return e
}

func (err *ErrorWithInternalCode) Error() string {
	return err.Err.Error()
}

func (err *ErrorWithInternalCode) Wrappee() error {
	return err.Err
}
