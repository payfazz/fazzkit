package http_test

import (
	"context"
	nethttp "net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/payfazz/fazzkit/server/http"
)

type Foo struct {
	A int        `httpquery:"a"`
	B *int       `httpquery:"b"`
	C *string    `httpquery:"c"`
	D *uuid.UUID `httpquery:"d"`
}

type Bar struct {
	A int        `scheme:"a"`
	B *int       `scheme:"b"`
	C *string    `scheme:"c"`
	D *uuid.UUID `scheme:"d"`
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
	q.Add("d", "39e0233b-9ffd-4718-abe2-17b6d6589ef2")

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

	id, _ := uuid.Parse("39e0233b-9ffd-4718-abe2-17b6d6589ef2")
	if *foo.D != id {
		t.Fatal("uuid expected")
	}
}

func TestDecodeURL(t *testing.T) {

	form := url.Values{}

	form.Add("a", "3")
	form.Add("b", "4")
	form.Add("c", "hello")
	form.Add("d", "39e0233b-9ffd-4718-abe2-17b6d6589ef2")
	request, err := nethttp.NewRequest(nethttp.MethodPost, "localhost", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Add("content-type", "application/x-www-form-urlencoded")

	decodeFunc := http.DecodeURLEncoded(&Bar{})

	result, err := decodeFunc(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}

	bar := result.(*Bar)
	if bar.A != 3 {
		t.Fatal("3 expected")
	}

	if *bar.B != 4 {
		t.Fatal("4 expected")
	}

	if *bar.C != "hello" {
		t.Fatal("hello expected")
	}

	id, _ := uuid.Parse("39e0233b-9ffd-4718-abe2-17b6d6589ef2")
	if *bar.D != id {
		t.Fatal("uuid expected")
	}
}
