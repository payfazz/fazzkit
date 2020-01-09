package error

type ErrorWithInternalCode struct {
	Code string
	Err  error
}

func NewErrorWithInternalCode(code string, err error) *ErrorWithInternalCode {
	return &ErrorWithInternalCode{
		Code: code,
		Err:  err,
	}
}

func (err *ErrorWithInternalCode) Error() string {
	return err.Err.Error()
}
