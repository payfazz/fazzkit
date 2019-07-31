package servererror

//ErrorWithStatusCode error with http status code
type ErrorWithStatusCode struct {
	Err        string
	StatusCode int
}

func (e *ErrorWithStatusCode) Error() string {
	return e.Err
}