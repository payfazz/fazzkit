package retry_test

import (
	"testing"
	"time"

	"github.com/nats-io/stan.go"
	"github.com/payfazz/fazzkit/event/stan/retry"
)

func TestDefaultOpt(t *testing.T) {
	opt := retry.NewOpt(retry.Opt{})
	if opt.Async != false {
		t.Error("async must be false")
	}
	if opt.Delay != time.Second*5 {
		t.Error("delay must be 5s")
	}
	if opt.IncrementDelayFunc(opt.Delay) != time.Second*5 {
		t.Error("increment must be stagnan")
	}
	if opt.UpTo != -1 {
		t.Error("default up to must be -1")
	}
	if opt.Logger == nil {
		t.Error("nil logger")
	}
}

var i = 0

func panicHandler(msg *stan.Msg) {
	i += 1
	panic("test")
}

func forcePanicHandler(msg *stan.Msg) {
	panic(retry.ForcePanic("force"))
}

func TestRetry(t *testing.T) {
	handler := retry.Retry(panicHandler, retry.NewOpt(retry.Opt{
		UpTo:  2,
		Delay: time.Second * 1,
	}))
	handler(&stan.Msg{})
	if i != 3 {
		t.Error("must be 3")
	}
}

var j = 0

func stopHandler(msg *stan.Msg) {
	j += 1
	if j == 2 {
		panic(retry.StopRetry("stop"))
	}

	panic("normal panic")
}

func TestStopRetry(t *testing.T) {
	handler := retry.Retry(stopHandler, retry.NewOpt(retry.Opt{
		UpTo:  -1,
		Delay: time.Second * 1,
	}))
	handler(&stan.Msg{})
	if j != 2 {
		t.Log(j)
		t.Error("must be 2")
	}
}

func forceHandler(msg *stan.Msg) {
	panic(retry.ForcePanic("force"))
}

func TestForcePanic(t *testing.T) {
	handler := retry.Retry(forceHandler, retry.NewOpt(retry.Opt{
		UpTo:  -1,
		Delay: time.Second * 1,
	}))

	defer func() {
		var err interface{}
		if err = recover(); err == nil {
			t.Error("must be panic")
		}
	}()

	handler(&stan.Msg{})
}

func TestASync(t *testing.T) {
	start := time.Now()
	defer func() {
		if time.Since(start) > time.Second*2 {
			t.Error("maybe not async")
		}
	}()

	handler := retry.Retry(func(msg *stan.Msg) { panic("a") }, retry.NewOpt(retry.Opt{
		UpTo:  2,
		Delay: time.Second * 1,
		Async: true,
	}))
	handler(&stan.Msg{})
	handler(&stan.Msg{})
}

var k = 0

func noRetry(msg *stan.Msg) {
	k += 1
}

func TestNoRetry(t *testing.T) {
	handler := retry.Retry(noRetry, retry.NewOpt(retry.Opt{
		UpTo:  2,
		Delay: time.Second * 1,
		Async: true,
	}))
	handler(&stan.Msg{})
	if k > 1 {
		t.Error("retried")
	}
}
