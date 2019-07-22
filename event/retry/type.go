package retry

import (
	"os"
	"time"

	"github.com/go-kit/kit/log"
)

//IncrementDelayFunc delay manipulation every retry
type IncrementDelayFunc func(time.Duration) time.Duration

//Opt retry option
type Opt struct {
	attempt            int
	UpTo               int
	Delay              time.Duration
	IncrementDelayFunc IncrementDelayFunc
	Logger             *log.Logger
}

//NewOpt construct new option
func NewOpt(opt Opt) *Opt {
	newOpt := &Opt{
		UpTo:               opt.UpTo,
		Delay:              opt.Delay,
		IncrementDelayFunc: opt.IncrementDelayFunc,
		Logger:             opt.Logger,
		attempt:            opt.attempt,
	}

	if newOpt.Logger == nil {
		newOpt.Logger = initDefaultLogger()
	}

	if newOpt.IncrementDelayFunc == nil {
		newOpt.IncrementDelayFunc = defaultIncrementDelayFunc
	}

	return newOpt
}

func defaultIncrementDelayFunc(delay time.Duration) time.Duration {
	return delay
}

func initDefaultLogger() *log.Logger {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	return &logger
}
