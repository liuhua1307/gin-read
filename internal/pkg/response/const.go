package response

import (
	"github.com/liuhua1307/gin-read/internal/consts"
)

const (
	UnKnownMsg = "unKnown"
	SuccessMsg = "success"
)

var (
	InvalidTokenError = newError(consts.Unauthorized, "invalid token error")
	ServerError       = newError(consts.ServerError, "internal server error")
	FormError         = newError(consts.RequestError, "form error")
)

var (
	SuccessWithoutDataResult = newSuccessResult(nil)
	UnknownErrorResult       = newErrorResult(consts.UnknownError, UnKnownMsg)
)
