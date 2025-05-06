package hest

import (
	"github.com/gin-gonic/gin"
	"github.com/hootuu/hyle/data/hcast"
	"github.com/hootuu/hyle/herr"
	"github.com/hootuu/hyle/hio"
)

func HelixHandle[REQ any, RESP any](ctx *gin.Context, callback func(req *REQ) (*RESP, *herr.Error)) {
	ctx.Request.Header.Get()
	//switch ctx.Request.Method {
	//case http.MethodPost
	//}
	//bodyBytes, nErr := ctx.GetRawData()
	//if nErr != nil {
	//	gLogger.Error("get request body data failed", zap.Error(nErr))
	//	ctx.JSON(
	//		http.StatusOK,
	//		rest.FailResponse[any](idx.New(), errors.System("get request body data failed")),
	//	)
	//	return
	//}
	//req, err := rest.UnmarshalRequest[REQ](bodyBytes)
	//if err != nil {
	//	gLogger.Error("unmarshal request failed", zap.Error(err.Native()))
	//	ctx.JSON(
	//		http.StatusOK,
	//		rest.FailResponse[RESP](idx.New(), errors.System("unmarshal request failed")),
	//	)
	//	return
	//}
	//bReqVerified := true
	//err = Guard(req.GuardID, func(pubKey []byte) {
	//	innerErr := req.Verify(pubKey)
	//	if innerErr != nil {
	//		gLogger.Error("guard.Verify failed", zap.String("req", req.JSON()), zap.Error(innerErr))
	//		bReqVerified = false
	//	}
	//})
	//if err != nil {
	//	gLogger.Error("guard failed", zap.Error(err.Native()))
	//	ctx.JSON(
	//		http.StatusOK,
	//		rest.FailResponse[RESP](req.ID, errors.System("guard failed")),
	//	)
	//	return
	//}
	//if !bReqVerified {
	//	ctx.JSON(
	//		http.StatusOK,
	//		rest.FailResponse[RESP](req.ID, errors.System("invalid signature")),
	//	)
	//	return
	//}
	//data, err := callback(req.Data)
	//if err != nil {
	//	gLogger.Error("["+ctx.Request.Method+"]",
	//		zap.String("URL", ctx.Request.URL.String()),
	//		zap.Any("req", req),
	//		zap.Error(err))
	//	ctx.JSON(http.StatusOK, rest.FailResponse[RESP](req.ID, err))
	//	return
	//}
	//if gLogger.Level() <= zapcore.InfoLevel {
	//	gLogger.Info("["+ctx.Request.Method+"]",
	//		zap.String("URL", ctx.Request.URL.String()),
	//		zap.Any("req", req),
	//		zap.Any("data", data),
	//	)
	//}
	//ctx.JSON(http.StatusOK, rest.NewResponse[RESP](req.ID, data))
}

func reqFromHeader[T any](ctx *gin.Context) (*hio.Request[T], *herr.Error) {
	var req hio.Request[T]
	req.ReqID = ctx.Request.Header.Get(hio.HttpHeaderReqID)
	req.TokenID = ctx.Request.Header.Get(hio.HttpHeaderTokenID)
	req.Timestamp = hcast.ToInt64(ctx.Request.Header.Get(hio.HttpHeaderTimestamp))
	req.Nonce = hcast.ToInt64(ctx.Request.Header.Get(hio.HttpHeaderNonce))
	req.Signature = ctx.Request.Header.Get(hio.HttpHeaderSignature)
	return &req, nil
}

func handleGet[REQ any, RESP any](ctx *gin.Context) (*RESP, *herr.Error) {

}
