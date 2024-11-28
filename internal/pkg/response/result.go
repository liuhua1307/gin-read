package response

import (
	"github.com/liuhua1307/gin-read/internal/consts"
)

type result struct {
	Code consts.ResponseCode `json:"status"`
	Msg  string              `json:"message"`
	Data any                 `json:"data"`
}

func newSuccessResult(data any) *result {
	return &result{
		Code: consts.Success,
		Msg:  SuccessMsg,
		Data: data,
	}
}

func newErrorResult(code consts.ResponseCode, msg string) *result {
	return &result{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}
