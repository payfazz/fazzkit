package retry

import (
	"fmt"
	"time"

	stan "github.com/nats-io/stan.go"
)

//Retry run time.AfterFunc if panic happened
func Retry(handler stan.MsgHandler, opt *Opt) stan.MsgHandler {
	return func(msg *stan.Msg) {
		defer func(opt *Opt) {
			if err := recover(); err == nil {
				return
			}

			if opt.UpTo == opt.attempt {
				return
			}

			go time.AfterFunc(opt.Delay, func() {
				fmt.Println(opt)
				fmt.Printf("retry %d/%d\n", opt.attempt, opt.UpTo)
				Retry(handler, opt)(msg)
			})

			opt.attempt = opt.attempt + 1
			opt.Delay = opt.IncrementDelayFunc(opt.Delay)
		}(NewOpt(*opt))

		handler(msg)
	}
}
