package response

import (
	"errors"

	"github.com/liuhua1307/gin-read/internal/consts"
)

var _ error = &responseError{}

type responseError struct {
	Msg  string
	Code consts.ResponseCode
}

func (r *responseError) Error() string {
	return r.Msg
}

func newError(code consts.ResponseCode, msg string) error {
	if _, ok := consts.ErrorCodeHTTPStatus[code]; !ok {
		code = consts.UnknownError
	}
	if len(msg) == 0 {
		msg = UnKnownMsg
	}
	return &responseError{
		Code: code,
		Msg:  msg,
	}
}

func errorHTTPStatus(err error) int {
	var (
		e          *responseError
		ok         bool
		httpStatus int
	)

	if ok = errors.As(err, &e); !ok {
		return consts.ErrorCodeHTTPStatus[consts.UnknownError]
	}
	if httpStatus, ok = consts.ErrorCodeHTTPStatus[e.Code]; !ok {
		return consts.ErrorCodeHTTPStatus[consts.UnknownError]
	}
	return httpStatus
}

func errorResult(err error) *result {
	var (
		e  *responseError
		ok bool
	)
	if ok = errors.As(err, &e); !ok {
		return UnknownErrorResult
	}
	return newErrorResult(e.Code, e.Msg)
}
