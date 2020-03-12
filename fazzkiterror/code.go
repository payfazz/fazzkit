package fazzkiterror

var DefaultInternalCode = "-1"

func GetInternalCode(err error) string {
	if e, ok := err.(*ErrorWithInternalCode); ok {
		return e.Code
	}

	if e, ok := err.(Wrapper); ok {
		return GetInternalCode(e.Wrappee())
	}

	return DefaultInternalCode
}

func GetDomainError(err error) string {
	if e, ok := err.(Wrapper); ok {
		return GetDomainError(e.Wrappee())
	}

	return err.Error()
}
