package exactlyonce

import (
	"fmt"

	stan "github.com/nats-io/stan.go"
	"github.com/payfazz/fazzkit/event/stan/message"
)

//ExactlyOnce wrap stan.MsgHandler with exactly once procedure
func ExactlyOnce(handler stan.MsgHandler, opt *Opt) stan.MsgHandler {
	if opt.Repository == nil {
		opt.Repository = message.NewBoltRepo(opt.DbName)
	}

	return func(msg *stan.Msg) {
		repo := opt.Repository
		seq := fmt.Sprintf("%s_%d", msg.Subject, msg.Sequence)
		seqb := []byte(seq)

		loadedMsg, _ := repo.Load(seqb)
		if loadedMsg != nil {
			return
		}

		handler(msg)

		repo.Save(message.Message{
			ID:   seqb,
			Data: msg.Data,
		})
	}
}
