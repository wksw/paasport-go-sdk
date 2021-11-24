package paasport

// Error Define
type Error struct {
	Code    int32  `json:"code"`
	ErrCode int32  `json:"err_code"`
	Message string `json:"message"`
	RequestMeta
}

// OneboxError onebox error define
type OneboxError struct {
	Code    int32  `json:"code"`
	ErrCode int32  `json:"errCode"`
	Message string `json:"errMessage"`
	Prompt  string `json:"prompt"`
	Data    []byte `json:"data"`
	Status  int32  `json:"status"`
	Success bool   `json:"success"`
	RequestMeta
}

// RequestMeta the common response data
type RequestMeta struct {
	RequestId     int64  `json:"request_id,string"`
	RequestMethod string `json:"request_method"`
}
