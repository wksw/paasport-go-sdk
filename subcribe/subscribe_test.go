package subscribe

import (
	"testing"

	paasport "gitee.com/paasport/go-sdk"
	pb "gitee.com/paasport/protos-repo/account/subscribe"
)

func TestSubscribe(t *testing.T) {
	client, err := paasport.NewClient("", "", "http://127.0.0.1:9091")
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
	subscribeClient := &pb.SubscribeSDKClient{
		C: client,
	}
	resp, rerr := subscribeClient.Subscribe(&pb.SubscribeReq{})
	if err != nil {
		t.Log(rerr)
		t.FailNow()
	}
	t.Log(resp)
}
