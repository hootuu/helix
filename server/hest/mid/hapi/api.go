package hapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hootuu/helix/components/hguard"
	"github.com/hootuu/hyle/data/hcast"
	"github.com/hootuu/hyle/herr"
	"github.com/hootuu/hyle/hio"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"net/http"
)

func ApiHandle[REQ any, RESP any](handle func(req *REQ) (*RESP, *herr.Error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req *hio.Request[REQ]
		var resp *hio.Response[RESP]
		defer func() {
			if err := ctx.Err(); err != nil {
				prefix := fmt.Sprintf("[%s]%s", ctx.Request.Method, ctx.Request.URL.Path)
				hlog.Err(prefix, zap.Any("req", req), zap.Any("resp", resp), zap.Error(err))
			}
		}()

		switch ctx.Request.Method {
		case http.MethodPost:
		default:
			_ = ctx.Error(herr.Of(hio.ReqMethodMustBePost, "http method must be POST"))
			return
		}
		req = getReqFromHeader[REQ](ctx)
		if err := req.PreVerify(); err != nil {
			_ = ctx.Error(err)
			return
		}
		bodyBytes, nErr := ctx.GetRawData()
		if nErr != nil {
			hlog.Err("hapi.ApiHandle: ctx.GetRawData()", zap.Error(nErr))
			_ = ctx.Error(herr.Of(hio.ReqParseBodyDataErr, "parse raw data err:"+nErr.Error()))
			return
		}
		nErr = hguard.GuardVerify(req.TokenID, bodyBytes, req.Signature)
		if nErr != nil {
			_ = ctx.Error(herr.Of(hio.ReqApiGuardVerifyErr, "guard verify failed"))
			return
		}
		err := req.Unmarshal(bodyBytes)
		if err != nil {
			hlog.Err("hapi.ApiHandle: req.Unmarshal", zap.Error(err.Native()))
			_ = ctx.Error(err)
			return
		}
		respData, err := handle(req.Data)
		if err != nil {
			hlog.Err("hapi.ApiHandle: handle", zap.Error(err.Native()))
			_ = ctx.Error(err)
			return
		}
		resp = hio.NewResponse[RESP](req.ReqID, respData)
		ctx.JSON(http.StatusOK, resp)
	}
}

func getReqFromHeader[T any](ctx *gin.Context) *hio.Request[T] {
	var req hio.Request[T]
	req.ReqID = ctx.Request.Header.Get(hio.HttpHeaderReqID)
	req.TokenID = ctx.Request.Header.Get(hio.HttpHeaderTokenID)
	req.Timestamp = hcast.ToInt64(ctx.Request.Header.Get(hio.HttpHeaderTimestamp))
	req.Nonce = hcast.ToInt64(ctx.Request.Header.Get(hio.HttpHeaderNonce))
	req.Signature = ctx.Request.Header.Get(hio.HttpHeaderSignature)
	return &req
}
