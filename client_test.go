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
