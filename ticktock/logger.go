package ticktock

import (
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

var gLogger = hlog.GetLogger("ticktock")

type logger struct {
	code string
}

func newLogger(code string) *logger {
	return &logger{code}
}

func (l *logger) Debug(args ...interface{}) {
	gLogger.Debug(l.code+"-D", zap.Any("args", args))
}

func (l *logger) Info(args ...interface{}) {
	gLogger.Info(l.code+"-I", zap.Any("args", args))
}

func (l *logger) Warn(args ...interface{}) {
	gLogger.Warn(l.code+"-W", zap.Any("args", args))
}

func (l *logger) Error(args ...interface{}) {
	gLogger.Error(l.code+"-E", zap.Any("args", args))
}

func (l *logger) Fatal(args ...interface{}) {
	gLogger.Fatal(l.code+"-F", zap.Any("args", args))
}
