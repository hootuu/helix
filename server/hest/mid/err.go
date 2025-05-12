package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/hootuu/hyle/herr"
	"github.com/hootuu/hyle/hio"
	"net/http"
)

func ErrHandle(ctx *gin.Context) {
	ctx.Next()
	ginErr := ctx.Err()
	if err, ok := ginErr.(*herr.Error); ok {
		reqID := ctx.Request.Header.Get(hio.HttpHeaderReqID)
		if reqID == "" {
			reqID, _ = ctx.GetQuery(hio.HttpHeaderReqID)
		}
		ctx.JSON(
			http.StatusOK,
			hio.FailResponse[any](reqID, err),
		)
	}
}
