package log

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	GL "gorm.io/gorm/logger"
)

// GormZapLogger zap.Logger 的包装 使用 zap.Logger 实现 logger.Interface
type GormZapLogger struct {
	ZapLogger     *zap.Logger
	SlowThreshold time.Duration // 慢查询阈值
}

func NewGormZapLogger(zapLog *zap.Logger, slowThreshold time.Duration) *GormZapLogger {
	return &GormZapLogger{
		ZapLogger:     zapLog,
		SlowThreshold: slowThreshold,
	}
}

func (l *GormZapLogger) LogMode(_ GL.LogLevel) GL.Interface {
	return l
}

func (l *GormZapLogger) Info(_ context.Context, msg string, data ...interface{}) {
	l.ZapLogger.Info(fmt.Sprintf(msg, data...))
}

func (l *GormZapLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	l.ZapLogger.Warn(fmt.Sprintf(msg, data...))
}

func (l *GormZapLogger) Error(_ context.Context, msg string, data ...interface{}) {
	l.ZapLogger.Error(fmt.Sprintf(msg, data...))
}

func (l *GormZapLogger) Trace(_ context.Context, begin time.Time, f func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := f()

	fields := []zap.Field{
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
		zap.String("sql", sql),
	}

	if err != nil {
		l.ZapLogger.Error("SQL Error", append(fields, zap.Error(err))...)
		return
	}

	if elapsed > l.SlowThreshold && l.SlowThreshold != 0 {
		l.ZapLogger.Warn("Slow SQL", fields...)
		return
	}

	l.ZapLogger.Debug("SQL Trace", fields...)
}
