//go:build !prod

package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	cfg := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		FunctionKey:      "func",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.LowercaseLevelEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: " ",
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(cfg), zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	logger = zap.New(core)
}
