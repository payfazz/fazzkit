package retry

import (
	"os"
	"time"

	"github.com/go-kit/kit/log"
)

//IncrementDelayFunc modify delay duration every attempt
type IncrementDelayFunc func(time.Duration) time.Duration

//Opt retry options
type Opt struct {
	UpTo               int
	Delay              time.Duration
	IncrementDelayFunc IncrementDelayFunc
	Logger             *log.Logger
	Async              bool
	attempt            int
}

//NewOpt create retry options
func NewOpt(opt Opt) *Opt {
	newOpt := &Opt{
		UpTo:               opt.UpTo,
		Delay:              opt.Delay,
		IncrementDelayFunc: opt.IncrementDelayFunc,
		Logger:             opt.Logger,
		Async:              opt.Async,
		attempt:            opt.attempt,
	}

	if newOpt.Delay == 0 {
		newOpt.Delay = time.Second * 5
	}

	if newOpt.UpTo == 0 {
		newOpt.UpTo = -1
	}

	if newOpt.Logger == nil {
		newOpt.Logger = initLogger()
	}

	if newOpt.IncrementDelayFunc == nil {
		newOpt.IncrementDelayFunc = defaultIncrementDelayFunc
	}

	return newOpt
}

func defaultIncrementDelayFunc(delay time.Duration) time.Duration {
	return delay
}

func initLogger() *log.Logger {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	return &logger
}
