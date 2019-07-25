package exactlyonce_test

import (
	"testing"
	"time"

	"github.com/nats-io/stan.go"
	"github.com/payfazz/fazzkit/event/stan/exactlyonce"
)

var i = 0

func foo(msg *stan.Msg) {
	i += 1
}

func TestExactlyOnce(t *testing.T) {
	opt := &exactlyonce.Opt{DbName: time.Now().Format("20060102150405")}
	handler := exactlyonce.ExactlyOnce(foo, exactlyonce.NewOpt(*opt))
	handler(&stan.Msg{})
	handler(&stan.Msg{})
	handler(&stan.Msg{})
	if i > 1 {
		t.Error("must be 1")
	}
}
