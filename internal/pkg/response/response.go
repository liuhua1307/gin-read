package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error(ctx *gin.Context, err error) {
	if err == nil {
		SuccessWithoutData(ctx)
		return
	}
	response(ctx, errorResult(err), errorHTTPStatus(err))
}

func SuccessWithData(ctx *gin.Context, data any) {
	response(ctx, newSuccessResult(data), http.StatusOK)
}

func SuccessWithoutData(ctx *gin.Context) {
	response(ctx, SuccessWithoutDataResult, http.StatusOK)
}

func response(ctx *gin.Context, data any, httpStatus int) {
	ctx.JSON(httpStatus, data)
	ctx.AbortWithStatus(httpStatus)
}
