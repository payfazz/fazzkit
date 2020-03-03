package fazzkiterror

type Wrapper interface {
	error
	Wrappee() error
}
