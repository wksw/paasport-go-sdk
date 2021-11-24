/*
	订阅相关接口
*/

package paasport

import (
	"net/http"

	pb "gitee.com/paasport/protos-repo/account/subscribe"
)

// Subscribe 订阅
func (c Client) Subscribe(in *pb.SubscribeReq) (out *pb.SubscribeResp, err *Error) {
	err = c.Do(http.MethodPost, "/subscribe", in, out)
	return out, err
}
