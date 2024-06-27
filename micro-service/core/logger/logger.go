package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugarLogger *zap.SugaredLogger

// 初始化日志配置
func NewLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 时间格式
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 级别格式

	config.OutputPaths = append(config.OutputPaths, "./logs")

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()
	sugarLogger = logger.Sugar()
}

// 获取全局日志对象
func GetLogger() *zap.SugaredLogger {
	if sugarLogger == nil {
		NewLogger()
	}

	return sugarLogger
}

// Info 日志
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Warn 日志
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Error 日志
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Debug 日志
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Panic 日志
func Panic(args ...interface{}) {
	GetLogger().Panic(args...)
}

// Infof 格式化日志
func Infof(template string, args ...interface{}) {
	GetLogger().Infof(template, args...)
}

// Warnf 格式化日志
func Warnf(template string, args ...interface{}) {
	GetLogger().Warnf(template, args...)
}

// Errorf 格式化日志
func Errorf(template string, args ...interface{}) {
	GetLogger().Errorf(template, args...)
}

// Debugf 格式化日志
func Debugf(template string, args ...interface{}) {
	GetLogger().Debugf(template, args...)
}

// Panicf 格式化日志
func Panicf(template string, args ...interface{}) {
	GetLogger().Panicf(template, args...)
}
