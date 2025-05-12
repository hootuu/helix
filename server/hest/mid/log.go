package mid

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hootuu/hyle/hio"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"strings"
	"time"
)

type LoggerMid struct {
	logger *zap.Logger
}

func NewLoggerMid(code string) *LoggerMid {
	return &LoggerMid{
		logger: hlog.GetLogger(strings.ToLower(code)),
	}
}

func (mid *LoggerMid) Handle(c *gin.Context) {
	startTime := time.Now()
	method := c.Request.Method
	path := c.Request.URL.Path
	reqID := getReqID(c)
	prefix := fmt.Sprintf("[%s]%s", method, path)

	mid.logger.Info(prefix, zap.String("reqID", reqID))

	c.Next()

	arr := []zap.Field{
		zap.String("reqID", reqID),
		zap.Int("s", c.Writer.Status()),
		zap.Duration("e", time.Since(startTime)),
	}
	if err := c.Err(); err != nil {
		arr = append(arr, zap.Error(err))
	}
	mid.logger.Info(prefix, arr...)
}

func getReqID(ctx *gin.Context) string {
	reqID := ctx.Request.Header.Get(hio.HttpHeaderReqID)
	if reqID == "" {
		reqID, _ = ctx.GetQuery(hio.HttpHeaderReqID)
	}
	return reqID
}
