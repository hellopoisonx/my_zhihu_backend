//go:build prod

package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	file, err := os.OpenFile("log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	cfg := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		FunctionKey:      zapcore.OmitKey,
		MessageKey:       "msg",
		StacktraceKey:    zapcore.OmitKey,
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.LowercaseLevelEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: " ",
	}
	syncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(file))
	core := zapcore.NewCore(zapcore.NewJSONEncoder(cfg), syncer, zapcore.InfoLevel)
	logger = zap.New(core)
}
