package logger

import "go.uber.org/zap"

var L LoggerInterface

type LoggerInterface interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

type WrappedLogger struct {
	Logger LoggerInterface
}

func (w WrappedLogger) Debug(msg string, fields ...zap.Field) {
	w.Logger.Debug(msg, insertField(fields)...)
}

func (w WrappedLogger) Info(msg string, fields ...zap.Field) {
	w.Logger.Info(msg, insertField(fields)...)
}

func (w WrappedLogger) Warn(msg string, fields ...zap.Field) {
	w.Logger.Warn(msg, insertField(fields)...)
}

func (w WrappedLogger) Error(msg string, fields ...zap.Field) {
	w.Logger.Error(msg, insertField(fields)...)
}

func insertField(fields []zap.Field) []zap.Field {
	ls := []zap.Field{zap.String("owner", "[map-chan]")}
	ls = append(ls, fields...)
	return ls
}
