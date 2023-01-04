package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// use this same interface throughout project
type Logger interface {
	Debugf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
}


func NewLogger() *zap.SugaredLogger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "/var/log/rpipoller/rpipoller.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)

	return zap.New(core).Sugar()
}