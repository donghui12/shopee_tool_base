// logger/gorm_logger.go
package logger

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

// GormLogger 实现了 gorm.Logger.Interface
type GormLogger struct {
	ZapLogger                 *zap.Logger
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogLevel                  logger.LogLevel
}

// NewGormLogger 创建一个新的 GORM Logger 实例
func NewGormLogger(slowThreshold time.Duration) *GormLogger {
	return &GormLogger{
		ZapLogger:                 log, // 使用你现有的全局 log 实例
		SlowThreshold:             slowThreshold,
		IgnoreRecordNotFoundError: true,        // 通常忽略"记录未找到"错误
		LogLevel:                  logger.Warn, // 默认日志级别
	}
}

// LogMode 实现 logger.Interface 的 LogMode 方法
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 实现 logger.Interface 的 Info 方法
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.ZapLogger.Info(fmt.Sprintf(msg, data...))
	}
}

// Warn 实现 logger.Interface 的 Warn 方法
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.ZapLogger.Warn(fmt.Sprintf(msg, data...))
	}
}

// Error 实现 logger.Interface 的 Error 方法
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.ZapLogger.Error(fmt.Sprintf(msg, data...))
	}
}

// Trace 实现 logger.Interface 的 Trace 方法 - 这是最关键的方法！
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 构建日志字段
	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	if err != nil && !(l.IgnoreRecordNotFoundError && IsRecordNotFoundError(err)) {
		fields = append(fields, zap.Error(err))
		l.ZapLogger.Error("SQL Execution Error", fields...)
		return
	}

	// 慢查询日志
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.ZapLogger.Warn(fmt.Sprintf("SLOW SQL (> %v)", l.SlowThreshold), fields...)
		return
	}

	// 普通 SQL 日志（根据日志级别决定是否记录）
	if l.LogLevel >= logger.Info {
		l.ZapLogger.Info("SQL Query", fields...)
	}
}

// IsRecordNotFoundError 检查是否为"记录未找到"错误
// 这是一个简单的实现，你可能需要根据实际情况调整
func IsRecordNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	// GORM 的 RecordNotFound 错误通常包含 "record not found"
	return fmt.Sprintf("%s", err) == "record not found"
}
