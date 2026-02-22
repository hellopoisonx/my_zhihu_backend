package log

import (
	"go.uber.org/zap"
)

func L() *zap.Logger {
	if logger == nil {
		panic("logger not initialized")
	}
	return logger
}
