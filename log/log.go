package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.SugaredLogger
)

func init() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stdout"}

	cfg.Encoding = "console"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var lvl zapcore.Level
	lvl.Set("debug")
	cfg.Level = zap.NewAtomicLevelAt(lvl)

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	log = l.Sugar()
}

func G() *zap.SugaredLogger {
	if log == nil {
		return zap.L().Sugar()
	}
	return log
}
