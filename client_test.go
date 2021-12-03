package paasport

import (
	"fmt"
	"net/http"
	"testing"
)

type health struct {
	Status          int32           `protobuf:"varint,1,opt,name=status,proto3,enum=grpc.health.v1.HealthCheckResponse_ServingStatus" json:"status" form:"status"`
	ServiceName     string          `protobuf:"bytes,2,opt,name=service_name,json=serviceName,proto3" json:"service_name" form:"service_name"`
	AppName         string          `protobuf:"bytes,3,opt,name=app_name,json=appName,proto3" json:"app_name" form:"app_name"`
	TenantName      string          `protobuf:"bytes,4,opt,name=tenant_name,json=tenantName,proto3" json:"tenant_name" form:"tenant_name"`
	Version         string          `protobuf:"bytes,5,opt,name=version,proto3" json:"version" form:"version"`
	BuildNum        string          `protobuf:"bytes,6,opt,name=build_num,json=buildNum,proto3" json:"build_num" form:"build_num"`
	Hostname        string          `protobuf:"bytes,7,opt,name=hostname,proto3" json:"hostname" form:"hostname"`
	Now             string          `protobuf:"bytes,8,opt,name=now,proto3" json:"now" form:"now"`
	DependiceHealth map[string]bool `protobuf:"bytes,9,rep,name=dependice_health,json=dependiceHealth,proto3" json:"dependice_health" form:"dependice_health" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func TestClient(t *testing.T) {
	client, err := NewClient("", "", "http://127.0.0.1:9091")
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
	var resp health
	if err := client.Do(http.MethodGet, "/", nil, &resp); err != nil {
		t.Logf("%+v", err)
		t.FailNow()
	}
	t.Logf("%+v", resp)
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
