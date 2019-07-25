package message_test

import (
	"testing"

	"github.com/payfazz/fazzkit/event/stan/message"
)

func TestBolt(t *testing.T) {
	repo := message.NewBoltRepo("test")
	err := repo.Save(message.Message{
		ID:   []byte("test-id"),
		Data: []byte("test-data"),
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	message, err := repo.Load([]byte("test-id"))
	if err != nil {
		t.Errorf(err.Error())
	}

	if string(message.Data) != "test-data" {
		t.Errorf("wrong data")
	}
}

func TestInvalidLoad(t *testing.T) {
	repo := message.NewBoltRepo("test")

	msg, _ := repo.Load([]byte("invalid_id"))
	if msg != nil {
		t.Errorf("expect nil return %v", msg)
	}
}

func TestFlyweight(t *testing.T) {
	message.NewBoltRepo("foo")
	message.NewBoltRepo("foo")
	message.NewBoltRepo("foo")
	message.NewBoltRepo("foo")
	message.NewBoltRepo("foo")

	if message.TotalBoltRepo() >= 5 {
		t.Error("cache not working")
	}
}

func TestClose(t *testing.T) {
	repo := message.NewBoltRepo("1")
	repo.Close()

	err := repo.Save(message.Message{
		ID:   []byte("test-id"),
		Data: []byte("test-data"),
	})

	if err == nil {
		t.Error("db closed error is expected")
	}

	_, err = repo.Load([]byte("id"))
	if err == nil {
		t.Error("db closed error is expected")
	}

	repo = message.NewBoltRepo("1")
	err = repo.Save(message.Message{
		ID:   []byte("test-id"),
		Data: []byte("test-data"),
	})

	if err != nil {
		t.Error("save failed")
	}
}
