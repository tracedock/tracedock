package logger

import (
	"go.uber.org/zap"
)

var zaplog *zap.Logger

func init() {
	zaplog = zap.Must(zap.NewProduction())
}

func Debug(msg string) {
	zaplog.Debug(msg)
}

func Info(msg string) {
	zaplog.Info(msg)
}

func Warn(msg string) {
	zaplog.Warn(msg)
}

func Error(msg string) {
	zaplog.Error(msg)
}

func Fatal(msg string) {
	zaplog.Fatal(msg)
}

func Sync() {
	zaplog.Sync()
}
