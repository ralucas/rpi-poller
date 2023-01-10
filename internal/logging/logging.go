package logging

import (
	"io"
	"os"

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

type LoggerOutput string

const (
	FileOutput   LoggerOutput = "file"
	StdoutOutput LoggerOutput = "stdout"
)

type LoggerConfig struct {
	Output LoggerOutput
	Level  zapcore.Level
}

func Writer(o LoggerOutput) io.Writer {
	switch o {
	case FileOutput:
		return &lumberjack.Logger{
			Filename:   "/var/log/rpipoller/rpipoller.log",
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		}
	case StdoutOutput:
		return os.Stdout
	}

	return nil
}

func NewLogger(conf LoggerConfig) *zap.SugaredLogger {
	w := zapcore.AddSync(Writer(conf.Output))

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		w,
		conf.Level,
	)

	return zap.New(core).Sugar()
}
