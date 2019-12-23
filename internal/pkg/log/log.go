package log

import (
	"errors"

	"github.com/albertwidi/go-project-example/internal/pkg/log/logger"
	"github.com/albertwidi/go-project-example/internal/pkg/log/logger/std"
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
	_debugLogger logger.Logger
	_infoLogger  logger.Logger
	_warnLogger  logger.Logger
	_errorLogger logger.Logger
	_fatalLogger logger.Logger

	errInvalidLevel  = errors.New("log: invalid log level")
	errInvalidLogger = errors.New("log: invalid logger")
)

func init() {
	backend, err := std.New(nil)
	if err != nil {
		return
	}
	SetLogger(backend)
}

// SetLogger to set default logger backend
func SetLogger(backend logger.Logger) {
	_debugLogger = backend
	_infoLogger = backend
	_warnLogger = backend
	_errorLogger = backend
	_fatalLogger = backend
}

// SetConfig to the current logger
func SetConfig(config *logger.Config) error {
	if err := _debugLogger.SetConfig(config); err != nil {
		return err
	}

	if err := _infoLogger.SetConfig(config); err != nil {
		return err
	}

	if err := _warnLogger.SetConfig(config); err != nil {
		return err
	}

	if err := _errorLogger.SetConfig(config); err != nil {
		return err
	}

	if err := _fatalLogger.SetConfig(config); err != nil {
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

// setLevel function set the log level to the desired level for defaultLogger and _debugLogger
// _debugLogger level can go to any level, but not with defaultLogger
// this to make sure _debugLogger to be disabled when level is > debug
// and defaultLogger to not overlap with _debugLogger
func setLevel(level logger.Level) {
	_debugLogger.SetLevel(level)
	_infoLogger.SetLevel(level)
	_warnLogger.SetLevel(level)
	_errorLogger.SetLevel(level)
	_fatalLogger.SetLevel(level)
}

// Debug function
func Debug(args ...interface{}) {
	_debugLogger.Debug(args...)
}

// Debugf function
func Debugf(format string, v ...interface{}) {
	_debugLogger.Debugf(format, v...)
}

// Debugw function
func Debugw(msg string, keyValues logger.KV) {
	_debugLogger.Debugw(msg, keyValues)
}

// Print function
func Print(v ...interface{}) {
	_infoLogger.Info(v...)
}

// Println function
func Println(v ...interface{}) {
	_infoLogger.Info(v...)
}

// Printf function
func Printf(format string, v ...interface{}) {
	_infoLogger.Infof(format, v...)
}

// Info function
func Info(args ...interface{}) {
	_infoLogger.Info(args...)
}

// Infof function
func Infof(format string, v ...interface{}) {
	_infoLogger.Infof(format, v...)
}

// Infow function
func Infow(msg string, keyValues logger.KV) {
	_infoLogger.Infow(msg, keyValues)
}

// Warn function
func Warn(args ...interface{}) {
	_warnLogger.Warn(args...)
}

// Warnf function
func Warnf(format string, v ...interface{}) {
	_warnLogger.Warnf(format, v...)
}

// Warnw function
func Warnw(msg string, keyValues logger.KV) {
	_warnLogger.Warnw(msg, keyValues)
}

// Error function
func Error(args ...interface{}) {
	_errorLogger.Error(args...)
}

// Errorf function
func Errorf(format string, v ...interface{}) {
	_errorLogger.Errorf(format, v...)
}

// Errorw function
func Errorw(msg string, keyValues logger.KV) {
	_errorLogger.Errorw(msg, keyValues)
}

// Fatal function
func Fatal(args ...interface{}) {
	_fatalLogger.Fatal(args...)
}

// Fatalf function
func Fatalf(format string, v ...interface{}) {
	_fatalLogger.Fatalf(format, v...)
}

// Fatalw function
func Fatalw(msg string, keyValues logger.KV) {
	_fatalLogger.Fatalw(msg, keyValues)
}
