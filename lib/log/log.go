package log

import (
	"errors"

	"github.com/albertwidi/kothak/lib/log/logger"
)

// Config of log
type Config struct {
	Level string
	// LogFile for log to file
	// this is not needed by default
	// application is expected to run in containerized environment
	LogFile    string
	DebugFile  string
	TimeFormat string
	// set true to log line numbers
	// make sure you understand the overhead when use this
	Caller bool
	// set true to colorize log, only work in console
	UseColor bool
	// use json format
	UseJSON bool
}

var (
	debugLogger logger.Logger
	infoLogger  logger.Logger
	warnLogger  logger.Logger
	errorLogger logger.Logger
	fatalLogger logger.Logger

	errInvalidLevel  = errors.New("invalid log level")
	errInvalidLogger = errors.New("invalid logger")
)

// Init log package
func Init(backend logger.Logger) {
	debugLogger = backend
	infoLogger = backend
	warnLogger = backend
	errorLogger = backend
	fatalLogger = backend
}

// SetConfig to the current logger
func SetConfig(config *logger.Config) error {
	if err := debugLogger.SetConfig(config); err != nil {
		return err
	}

	if err := infoLogger.SetConfig(config); err != nil {
		return err
	}

	if err := warnLogger.SetConfig(config); err != nil {
		return err
	}

	if err := errorLogger.SetConfig(config); err != nil {
		return err
	}

	if err := fatalLogger.SetConfig(config); err != nil {
		return err
	}

	return nil
}

// SetLevel of log
func SetLevel(level logger.Level) {
	setLevel(level)
}

// SetLevelString to set log level using string
func SetLevelString(level string) {
	setLevel(logger.StringToLevel(level))
}

// setLevel function set the log level to the desired level for defaultLogger and debugLogger
// debugLogger level can go to any level, but not with defaultLogger
// this to make sure debugLogger to be disabled when level is > debug
// and defaultLogger to not overlap with debugLogger
func setLevel(level logger.Level) {
	debugLogger.SetLevel(level)
	infoLogger.SetLevel(level)
	warnLogger.SetLevel(level)
	errorLogger.SetLevel(level)
	fatalLogger.SetLevel(level)
}

// Debug function
func Debug(args ...interface{}) {
	debugLogger.Debug(args...)
}

// Debugf function
func Debugf(format string, v ...interface{}) {
	debugLogger.Debugf(format, v...)
}

// Debugw function
func Debugw(msg string, keyValues logger.KV) {
	debugLogger.Debugw(msg, keyValues)
}

// Print function
func Print(v ...interface{}) {
	infoLogger.Info(v...)
}

// Println function
func Println(v ...interface{}) {
	infoLogger.Info(v...)
}

// Printf function
func Printf(format string, v ...interface{}) {
	infoLogger.Infof(format, v...)
}

// Info function
func Info(args ...interface{}) {
	infoLogger.Info(args...)
}

// Infof function
func Infof(format string, v ...interface{}) {
	infoLogger.Infof(format, v...)
}

// Infow function
func Infow(msg string, keyValues logger.KV) {
	infoLogger.Infow(msg, keyValues)
}

// Warn function
func Warn(args ...interface{}) {
	warnLogger.Warn(args...)
}

// Warnf function
func Warnf(format string, v ...interface{}) {
	warnLogger.Warnf(format, v...)
}

// Warnw function
func Warnw(msg string, keyValues logger.KV) {
	warnLogger.Warnw(msg, keyValues)
}

// Error function
func Error(args ...interface{}) {
	errorLogger.Error(args...)
}

// Errorf function
func Errorf(format string, v ...interface{}) {
	errorLogger.Errorf(format, v...)
}

// Errorw function
func Errorw(msg string, keyValues logger.KV) {
	errorLogger.Errorw(msg, keyValues)
}

// Errors function to log errors package
func Errors(err error) {
	errorLogger.Errors(err)
}

// Fatal function
func Fatal(args ...interface{}) {
	fatalLogger.Fatal(args...)
}

// Fatalf function
func Fatalf(format string, v ...interface{}) {
	fatalLogger.Fatalf(format, v...)
}

// Fatalw function
func Fatalw(msg string, keyValues logger.KV) {
	fatalLogger.Fatalw(msg, keyValues)
}
