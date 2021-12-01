package paasport

import (
	"net/http"
	"testing"
)

func TestClient(t *testing.T) {
	client, err := NewClient("", "", "http://127.0.0.1:9091")
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
	if err := client.Do(http.MethodGet, "/", nil, nil); err != nil {
		t.Log(err.Message)
		t.FailNow()
	}
}

func TestRequest(t *testing.T) {
	client, err := NewClient("", "", "http://127.0.0.1:9091")
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
	var requests = []struct {
		method     string
		path       string
		body       interface{}
		expectPath string
	}{
		{method: http.MethodGet, path: "/a/b", body: struct {
			A string `json:"a"`
		}{A: "a"}, expectPath: "http://127.0.0.1:9091/a1/a/b?a=a"},
		{method: http.MethodPost, path: "/a/b", body: struct {
			A string `json:"a"`
		}{A: "a"}, expectPath: "http://127.0.0.1:9091/a1/a/b"},
		{method: http.MethodGet, path: "/a/{a}/b", body: struct {
			A string `json:"a"`
		}{A: "1"}, expectPath: "http://127.0.0.1:9091/a1/a/1/b?a=1"},
		{method: http.MethodPost, path: "/a/{a}/b", body: struct {
			A string `json:"a"`
		}{A: "1"}, expectPath: "http://127.0.0.1:9091/a1/a/1/b"},
	}
	for _, request := range requests {
		requestPath, _, err := client.request(request.method, request.path, request.body)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if requestPath != request.expectPath {
			t.Errorf("return '%s' want '%s'", requestPath, request.expectPath)
			t.FailNow()
		}
	}
}
