package error

type Wrapper interface {
	error
	Wrappee() error
}
