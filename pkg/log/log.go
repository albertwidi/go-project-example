package log

import (
	"github.com/albertwidi/kothak/pkg/log/logger"
	"go.uber.org/zap"
)

// level of log
const (
	DebugLevel = logger.DebugLevel
	InfoLevel  = logger.InfoLevel
	WarnLevel  = logger.WarnLevel
	ErrorLevel = logger.ErrorLevel
	FatalLevel = logger.FatalLevel
)

// Log level
const (
	DebugLevelString = logger.DebugLevelString
	InfoLevelString  = logger.InfoLevelString
	WarnLevelString  = logger.WarnLevelString
	ErrorLevelString = logger.ErrorLevelString
	FatalLevelString = logger.FatalLevelString
)

var defaultLogger *logger.Log

func init() {
	defaultLogger = logger.Default()
}

// SetOutputToFile function
func SetOutputToFile(filepath string) error {
	return defaultLogger.SetOutputToFile(filepath)
}

// SetLevel will set level to logger and create a new logger based on level
func SetLevel(l logger.Level) {
	defaultLogger.SetLevel(l)
}

// GetLevel return log level in string
func GetLevel() string {
	return defaultLogger.GetLevel()
}

// SetLevelString set level from string level
func SetLevelString(l string) {
	defaultLogger.SetLevelString(l)
}

// Debug log
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

// Debugf log
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Debugw log
func Debugw(msg string, keyAndValues ...interface{}) {
	defaultLogger.Debugw(msg, keyAndValues...)
}

// Print log
func Print(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Println log
func Println(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Printf log
func Printf(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Printw log
func Printw(msg string, keyAndValues ...interface{}) {
	defaultLogger.Infow(msg, keyAndValues...)
}

// Info log
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Infof log
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Infow log
func Infow(msg string, keyAndValues ...interface{}) {
	defaultLogger.Infow(msg, keyAndValues...)
}

// Warn log
func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Warnf log
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

// Warnw log
func Warnw(msg string, keyAndValues ...interface{}) {
	defaultLogger.Warnw(msg, keyAndValues...)
}

// Error log
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Errorf log
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// Errorw log
func Errorw(msg string, keyAndValues ...interface{}) {
	defaultLogger.Errorw(msg, keyAndValues...)
}

// Errors log log error detail from Errs
func Errors(err error) {
	defaultLogger.Errors(err)
}

// Fatal log
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

// Fatalf log
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

// Fatalw log
func Fatalw(format string, keyAndValues ...interface{}) {
	defaultLogger.Fatalw(format, keyAndValues...)
}

// With log
func With(args ...interface{}) *zap.SugaredLogger {
	return defaultLogger.With(args...)
}

// Config return logger config
func Config() logger.Config {
	return defaultLogger.Config()
}
