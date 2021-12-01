package paasport

import (
	"fmt"
	"net/http"
	"testing"
)

func TestClient(t *testing.T) {
	client, err := NewClient("", "", "http://127.0.0.1:9091")
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
	if err := client.Do(http.MethodPost, "/subscribe", nil, nil); err != nil {
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
		expectBody []byte
	}{
		{method: http.MethodGet, path: "/a/b", body: struct {
			A string `json:"a"`
		}{A: "a"}, expectPath: "http://127.0.0.1:9091/a1/a/b?a=a", expectBody: nil},
		{method: http.MethodPost, path: "/a/b", body: struct {
			A string `json:"a"`
		}{A: "a"}, expectPath: "http://127.0.0.1:9091/a1/a/b", expectBody: []byte(`{"a":"a"}`)},
		{method: http.MethodGet, path: "/a/{a}/b", body: struct {
			A string `json:"a"`
		}{A: "1"}, expectPath: "http://127.0.0.1:9091/a1/a/1/b?a=1", expectBody: nil},
		{method: http.MethodPost, path: "/a/{a}/b", body: struct {
			A int `json:"a"`
		}{A: 1}, expectPath: "http://127.0.0.1:9091/a1/a/1/b", expectBody: []byte(`{"a":1}`)},
	}
	for index, request := range requests {
		requestPath, requestBody, err := client.request(request.method, request.path, request.body)
		fmt.Println("----", string(requestBody))
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if requestPath != request.expectPath {
			t.Errorf("%d not expect request path return '%s' want '%s'", index, requestPath, request.expectPath)
			t.FailNow()
		}

		if string(requestBody) != string(request.expectBody) {
			t.Errorf("%d not expect request body return '%s' want '%s'", index, string(requestBody), string(request.expectBody))
			t.FailNow()
		}
	}
}
