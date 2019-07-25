package retry

import (
	"fmt"
	"time"

	stan "github.com/nats-io/stan.go"
)

//Retry wrap stan.MsgHandler with retry procedure
func Retry(handler stan.MsgHandler, opt *Opt) stan.MsgHandler {
	return func(msg *stan.Msg) {
		defer func(opt *Opt) {
			var err interface{}
			if err = recover(); err == nil {
				return
			}

			switch err.(type) {
			case *stopRetry:
				logger := *opt.Logger
				logger.Log("stop retry", err.(*stopRetry).Error())
				return
			case *forcePanic:
				logger := *opt.Logger
				logger.Log("force panic", err.(*forcePanic).Error())
				panic(err)
			}

			if opt.UpTo == opt.attempt {
				logger := *opt.Logger
				logger.Log("max_attempt_on", fmt.Sprintf("subject: %s sequence:%d", msg.Subject, msg.Sequence))
				return
			}

			opt.attempt = opt.attempt + 1

			if opt.Async {
				go time.AfterFunc(opt.Delay, func() {
					logRetry(opt)
					opt.Delay = opt.IncrementDelayFunc(opt.Delay)
					Retry(handler, opt)(msg)
				})
				return
			}

			time.Sleep(opt.Delay)
			logRetry(opt)
			opt.Delay = opt.IncrementDelayFunc(opt.Delay)
			Retry(handler, opt)(msg)

		}(NewOpt(*opt))
		handler(msg)
	}
}

func logRetry(opt *Opt) {
	logger := *opt.Logger

	var upTo string
	if opt.UpTo == -1 {
		upTo = "infinite"
	} else {
		upTo = fmt.Sprintf("%d", opt.UpTo)
	}
	logger.Log(
		"attempt", opt.attempt,
		"up_to", upTo,
		"delayed", opt.Delay,
		"next_attempt_delay", opt.IncrementDelayFunc(opt.Delay),
	)
}
