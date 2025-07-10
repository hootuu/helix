package htick

import (
	"fmt"
	"github.com/hootuu/hyle/hlog"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var gLogger = hlog.GetLogger("tick")
var gCron *cron.Cron

func init() {
	gCron = cron.New(cron.WithSeconds(), cron.WithLogger(&cronLogger{}))
	gCron.Start()
}

type cronLogger struct{}

func (c *cronLogger) Info(msg string, keysAndValues ...interface{}) {
	var d []zap.Field
	for idx, kv := range keysAndValues {
		d = append(d, zap.Any(fmt.Sprintf("%d", idx), kv))
	}
	gLogger.Debug(msg, d...)
}

func (c *cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	if keysAndValues == nil {
		gLogger.Error(msg, zap.Error(err))
		return
	}
	d := []zap.Field{zap.Error(err)}
	for idx, kv := range keysAndValues {
		d = append(d, zap.Any(fmt.Sprintf("%d", idx), kv))
	}
	gLogger.Error(msg, d...)
}
