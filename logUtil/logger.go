package logUtil

import (
	"fmt"
	"go.uber.org/zap"
)

type Logger struct {
	coreLoggger     *zap.Logger
	InnerModuleName string
}

func newLogger(logger *zap.Logger, moduleName string) *Logger {
	loggerWrapper := new(Logger)
	loggerWrapper.InnerModuleName = moduleName
	loggerWrapper.coreLoggger = logger
	return loggerWrapper
}

func (this *Logger) Info(msg string, fields ...zap.Field) {
	msg = fmt.Sprintf("[%s]  %s", this.InnerModuleName, msg)
	this.coreLoggger.Info(msg, fields...)
}
func (this *Logger) Warn(msg string, fields ...zap.Field) {
	msg = fmt.Sprintf("[%s]  %s", this.InnerModuleName, msg)
	this.coreLoggger.Warn(msg, fields...)
}
func (this *Logger) Error(msg string, fields ...zap.Field) {
	msg = fmt.Sprintf("[%s]  %s", this.InnerModuleName, msg)
	this.coreLoggger.Error(msg, fields...)
}
func (this *Logger) DPanic(msg string, fields ...zap.Field) {
	msg = fmt.Sprintf("[%s]  %s", this.InnerModuleName, msg)
	this.coreLoggger.DPanic(msg, fields...)
}
func (this *Logger) Panic(msg string, fields ...zap.Field) {
	msg = fmt.Sprintf("[%s]  %s", this.InnerModuleName, msg)
	this.coreLoggger.Panic(msg, fields...)
}
func (this *Logger) Fatal(msg string, fields ...zap.Field) {
	msg = fmt.Sprintf("[%s]  %s", this.InnerModuleName, msg)
	this.coreLoggger.Fatal(msg, fields...)
}
