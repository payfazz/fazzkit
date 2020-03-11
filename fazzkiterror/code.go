package fazzkiterror

var defaultInternalCode = "-1"

func SetDefaultInternalCode(code string) {
	defaultInternalCode = code
}

func GetInternalCode(err error) string {
	if e, ok := err.(*ErrorWithInternalCode); ok {
		return e.Code
	}

	if e, ok := err.(Wrapper); ok {
		return GetInternalCode(e.Wrappee())
	}

	return defaultInternalCode
}

func GetDomainError(err error) string {
	if e, ok := err.(Wrapper); ok {
		return GetDomainError(e.Wrappee())
	}

	return err.Error()
}
