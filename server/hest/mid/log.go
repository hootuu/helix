package mid

import (
	"fmt"
	"github.com/gin-gonic/gin"
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

func (mid *LoggerMid) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {

		startTime := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path

		c.Next()

		mid.logger.Info(fmt.Sprintf("[%s]%s", method, path),
			zap.Time("t", startTime),
			zap.Int("s", c.Writer.Status()),
			zap.Duration("e", time.Since(startTime)),
		)
	}
}
