package pblog

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func init() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	Logger = l.Sugar()
}

func Warnw(msg string, args ...interface{}) {
	Logger.Warnw(msg, args...)
}
