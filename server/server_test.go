package server_test

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/payfazz/fazzkit/server"
)

func Foo(ctx context.Context, request interface{}) (response interface{}, err error) {
	return nil, nil
}

func TestFoo(t *testing.T) {
	m := chi.NewRouter()
	m.Get("/foo", server.NewHTTPServer(Foo, server.HTTPOption{
		DecodeModel: nil,
		Logger:      nil,
	}).ServeHTTP)

	ts := httptest.NewServer(m)
	defer ts.Close()

	res, _ := testRequest(t, ts, "GET", "/foo", nil)
	if res.StatusCode != 204 {
		t.Fatal(res)
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}
