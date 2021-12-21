package paasport

import (
	"encoding/json"
	"testing"
)

func TestErrorSuccessMessage(t *testing.T) {
	value := `{
		"code": "123",
		"session_key": ""
	}`
	var respError Error
	if err := json.Unmarshal([]byte(value), &respError); err != nil {
		t.Fatal(err.Error())
		t.FailNow()
	}
	if respError.Code != 0 {
		t.FailNow()
	}
	t.Logf("%+v", respError)
}

func TestErrorFailMessage(t *testing.T) {
	value := `{
		"code": 1,
		"err_code": 123,
		"message": "error",
		"request_id": "123",
		"request_method": "123"
	}`
	var respError Error
	if err := json.Unmarshal([]byte(value), &respError); err != nil {
		t.Fatal(err.Error())
		t.FailNow()
	}
	if respError.Code == 0 {
		t.Logf("%+v", respError)
		t.FailNow()
	}
	t.Logf("%+v", respError)
}

func TestOneboxErrorSuccessMessage(t *testing.T) {
	value := `{
		"code": "123",
		"session_key": "",
		"data": {
			"a": "a"
		}
	}`
	var respError OneboxError
	if err := json.Unmarshal([]byte(value), &respError); err != nil {
		t.Fatal(err.Error())
		t.FailNow()
	}
	if respError.Code != 0 {
		t.FailNow()
	}
	t.Logf("%+v", respError)
}

func TestOneboxErrorFailMessage(t *testing.T) {
	value := `{
		"code": 1,
		"errCode": 123,
		"errMessage": "error",
		"prompt": "prompt",
		"request_id": "123",
		"request_method": "123"
	}`
	var respError OneboxError
	if err := json.Unmarshal([]byte(value), &respError); err != nil {
		t.Fatal(err.Error())
		t.FailNow()
	}
	if respError.Code == 0 {
		t.Logf("%+v", respError)
		t.FailNow()
	}
	t.Logf("%+v", respError)
}
