package retry

type stopRetry struct {
	err string
}

type forcePanic struct {
	err string
}

func (e *stopRetry) Error() string {
	return e.err
}

func (e *forcePanic) Error() string {
	return e.err
}

//StopRetry create error to stop retry procedure
func StopRetry(err string) error {
	return &stopRetry{err}
}

//ForcePanic create error to make panic unrecoverable in retry procedure
func ForcePanic(err string) error {
	return &forcePanic{err}
}
