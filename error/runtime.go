package error

type RuntimeError struct {
	RuntimeError error
	DomainError  error
}

func NewRuntimeError(domain, runtime error) Wrapper {
	return &RuntimeError{
		RuntimeError: runtime,
		DomainError:  domain,
	}
}

func (err *RuntimeError) Error() string {
	return err.RuntimeError.Error()
}


func (err *RuntimeError) Wrappee() error {
	return err.DomainError
}