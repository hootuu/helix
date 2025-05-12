package hest

import (
	"github.com/gin-gonic/gin"
	"github.com/hootuu/hyle/data/hcast"
	"github.com/hootuu/hyle/herr"
	"github.com/hootuu/hyle/hio"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"net/http"
)

func HelixHandle[REQ any, RESP any](ctx *gin.Context, callback func(req *REQ) (*RESP, *herr.Error)) *herr.Error {
	req, err := reqFromHeader[REQ](ctx)
	if err != nil {
		return err
	}
	switch ctx.Request.Method {
	case http.MethodPost:
	default:
		return herr.Of(hio.ReqMethodMustBePost, "http method must be POST")
	}
	bodyBytes, nErr := ctx.GetRawData()
	if nErr != nil {
		hlog.Err("hest.HelixHandle: ctx.GetRawData()", zap.Error(nErr))
		return herr.Of(hio.ReqParseBodyDataErr, "parse raw data err:"+nErr.Error())
	}
	err = req.Unmarshal(bodyBytes)
	if err != nil {
		hlog.Err("hest.HelixHandle: req.Unmarshal", zap.Error(err.Native()))
		return err
	}
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
	return nil
}

func reqFromHeader[T any](ctx *gin.Context) (*hio.Request[T], *herr.Error) {
	var req hio.Request[T]
	req.ReqID = ctx.Request.Header.Get(hio.HttpHeaderReqID)
	req.TokenID = ctx.Request.Header.Get(hio.HttpHeaderTokenID)
	req.Timestamp = hcast.ToInt64(ctx.Request.Header.Get(hio.HttpHeaderTimestamp))
	req.Nonce = hcast.ToInt64(ctx.Request.Header.Get(hio.HttpHeaderNonce))
	req.Signature = ctx.Request.Header.Get(hio.HttpHeaderSignature)
	if err := req.PreVerify(); err != nil {
		return nil, err
	}
	return &req, nil
}

func handleGet[REQ any, RESP any](ctx *gin.Context) (*RESP, *herr.Error) {

}
