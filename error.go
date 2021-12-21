package paasport

import "encoding/json"

// Error Define
type Error struct {
	Code       int32  `json:"code"`
	StatusCode int    `json:"status_code"`
	ErrCode    int32  `json:"err_code"`
	Message    string `json:"message"`
	RequestMeta
}

// ResponseError 错误返回
type ResponseError struct {
	ErrCode int32 `json:"err_code"`
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (e *Error) UnmarshalJSON(value []byte) error {
	var respError ResponseError
	if err := json.Unmarshal(value, &respError); err != nil {
		return nil
	}
	if respError.ErrCode != 0 {
		type alias Error
		aux := &struct {
			*alias
		}{
			alias: (*alias)(e),
		}
		if err := json.Unmarshal(value, &aux); err != nil {
			return err
		}
	}
	return nil
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

// ResponseOneboxError onebox错误返回
type ResponseOneboxError struct {
	ErrCode int32 `json:"errCode"`
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (e *OneboxError) UnmarshalJSON(value []byte) error {
	var respError ResponseOneboxError
	if err := json.Unmarshal(value, &respError); err != nil {
		return nil
	}
	if respError.ErrCode != 0 {
		type alias OneboxError
		aux := &struct {
			*alias
		}{
			alias: (*alias)(e),
		}
		if err := json.Unmarshal(value, &aux); err != nil {
			return err
		}
	}
	return nil
}

// RequestMeta the common response data
type RequestMeta struct {
	RequestId     int64  `json:"request_id,string"`
	RequestMethod string `json:"request_method"`
}
