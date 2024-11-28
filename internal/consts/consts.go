package consts

import "net/http"

type ResponseCode uint8

const (
	ServerError ResponseCode = iota
	Success
	UnknownError
	Forbidden
	Unauthorized
	RequestError
)

var ErrorCodeHTTPStatus = map[ResponseCode]int{
	UnknownError: http.StatusBadRequest,
	Success:      http.StatusOK,
	ServerError:  http.StatusInternalServerError,
	Forbidden:    http.StatusForbidden,
	Unauthorized: http.StatusUnauthorized,
	RequestError: http.StatusBadRequest,
}
