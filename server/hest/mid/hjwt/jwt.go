package hjwt

import (
	"github.com/gin-gonic/gin"
	"github.com/hootuu/hyle/data/hcast"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/herr"
	"github.com/hootuu/hyle/hio"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"net/http"
)

func JwtLogin[REQ any, RESP any](login func(req *REQ) (*RESP, *herr.Error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqID := getReqID(ctx)
		bodyBytes, nErr := ctx.GetRawData()
		if nErr != nil {
			hlog.Err("hjwt.JwtLogin: ctx.GetRawData()", zap.Error(nErr))
			err := herr.Of(hio.ReqParseBodyDataErr, "parse raw data err:"+nErr.Error())
			ctx.JSON(http.StatusOK, hio.FailResponse[any](reqID, err))
			return
		}
		req, nErr := hjson.FromBytes[REQ](bodyBytes)
		if nErr != nil {
			hlog.Err("hjwt.JwtLogin: hjson.FromBytes()", zap.Error(nErr))
			err := herr.Of(hio.ReqParseBodyDataErr, "parse raw data err:"+nErr.Error())
			ctx.JSON(http.StatusOK, hio.FailResponse[any](reqID, err))
			return
		}
		resp, err := login(req)
		if err != nil {
			ctx.JSON(http.StatusOK, hio.FailResponse[any](reqID, err))
			return
		}
		ctx.JSON(http.StatusOK, hio.NewResponse(reqID, resp))
	}
}

func JwtRefresh() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func JwtLogout(afterLogout func()) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenID := getTokenID(ctx)
		if tokenID == "" {
			err := herr.Of(hio.ReqRequireTokenID, "require h_token_id in header or query")
			ctx.JSON(http.StatusOK, hio.FailResponse(getReqID(ctx), err))
			return
		}
		//todo do logout
		afterLogout()
	}
}

func JwtHandle[REQ any, RESP any](handle func(req *REQ) (*RESP, *herr.Error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		switch ctx.Request.Method {
		case http.MethodPost:
		default:
			err := herr.Of(hio.ReqMethodMustBePost, "http method must be POST")
			ctx.JSON(http.StatusOK, hio.FailResponse[RESP](ctx.Request.Header.Get(hio.HttpHeaderReqID), err))
			return
		}
		req := getReqFromHeader[REQ](ctx)
		if err := req.PreVerify(); err != nil {
			hlog.Err("hjwt.JwtHandle: req.PreVerify", zap.Error(err))
			ctx.JSON(http.StatusOK, hio.FailResponse[RESP](req.ReqID, err))
			return
		}
		bodyBytes, nErr := ctx.GetRawData()
		if nErr != nil {
			hlog.Err("hjwt.JwtHandle: ctx.GetRawData()", zap.Error(nErr))
			err := herr.Of(hio.ReqParseBodyDataErr, "parse raw data err:"+nErr.Error())
			ctx.JSON(http.StatusOK, hio.FailResponse[RESP](req.ReqID, err))
			return
		}
		err := req.Unmarshal(bodyBytes)
		if err != nil {
			hlog.Err("hjwt.JwtHandle: req.Unmarshal", zap.Error(err.Native()))
			ctx.JSON(http.StatusOK, hio.FailResponse[RESP](req.ReqID, err))
			return
		}
		resp, err := handle(req.Data)
		if err != nil {
			hlog.Err("hjwt.JwtHandle: handle", zap.Error(err.Native()))
			ctx.JSON(http.StatusOK, hio.FailResponse[RESP](req.ReqID, err))
			return
		}
		ctx.JSON(http.StatusOK, hio.NewResponse[RESP](req.ReqID, resp))
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

func getReqID(ctx *gin.Context) string {
	reqID := ctx.Request.Header.Get(hio.HttpHeaderReqID)
	if reqID == "" {
		reqID, _ = ctx.GetQuery(hio.HttpHeaderReqID)
	}
	return reqID
}

func getTokenID(ctx *gin.Context) string {
	tokenID := ctx.Request.Header.Get(hio.HttpHeaderTokenID)
	if tokenID == "" {
		tokenID, _ = ctx.GetQuery(hio.HttpHeaderTokenID)
	}
	return tokenID
}
