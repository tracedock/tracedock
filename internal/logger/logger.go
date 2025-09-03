package logger

import (
	"go.uber.org/zap"
)

var zaplog *zap.Logger

func init() {
	zaplog = zap.Must(zap.NewProduction())
}

func Info(msg string) {
	zaplog.Info(msg)
}

func Infof(msg string) {
	zaplog.Info(msg)
}

func Debug(msg string) {
	zaplog.Debug(msg)
}

func Debugf(msg string) {
	zaplog.Debug(msg)
}

func Error(msg string) {
	zaplog.Error(msg)
}

func Errorf(msg string) {
	zaplog.Error(msg)
}

func Fatal(msg string) {
	zaplog.Fatal(msg)
}

func Fatalf(msg string) {
	zaplog.Fatal(msg)
}

func Sync() {
	zaplog.Sync()
}
