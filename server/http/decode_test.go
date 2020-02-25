package http_test

import (
	"context"
	"github.com/payfazz/fazzkit/server/http"
	nethttp "net/http"
	"testing"
)

type Foo struct {
	A int `httpquery:"a"`
	B *int `httpquery:"b"`
	C *string `httpquery:"c"`
}

func TestDecode(t *testing.T) {
	request, err := nethttp.NewRequest(nethttp.MethodGet, "localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := request.URL.Query()
	q.Add("a", "3")
	q.Add("b", "4")
	q.Add("c", "hello")
	request.URL.RawQuery = q.Encode()

	decodeFunc := http.Decode(&Foo{})

	result, err := decodeFunc(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}

	foo := result.(*Foo)
	if foo.A != 3 {
		t.Fatal("3 expected")
	}

	if *foo.B != 4 {
		t.Fatal("4 expected")
	}

	if *foo.C != "hello" {
		t.Fatal("hello expected")
	}
}
