package logger

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type (
	// KV is a type for logging with more information
	// this used by with function
	KV map[string]interface{}

	// Logger interface
	Logger interface {
		// this method is not concurrently safe to acess
		// preferable to create a new logger instead
		SetConfig(config *Config) error
		SetLevel(level Level) error
		Debug(args ...interface{})
		Debugf(format string, args ...interface{})
		Debugw(msg string, KV KV)
		Info(args ...interface{})
		Infof(format string, args ...interface{})
		Infow(msg string, KV KV)
		Warn(args ...interface{})
		Warnf(format string, args ...interface{})
		Warnw(msg string, KV KV)
		Error(args ...interface{})
		Errorf(format string, args ...interface{})
		Errorw(msg string, KV KV)
		Fatal(args ...interface{})
		Fatalf(format string, args ...interface{})
		Fatalw(msg string, KV KV)
	}

	// Level of log
	Level int

	// Config of logger
	Config struct {
		Level      Level
		LogFile    string
		TimeFormat string
		Caller     bool
		UseColor   bool
		UseJSON    bool
	}
)

// list of log level
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Log level
const (
	DebugLevelString = "debug"
	InfoLevelString  = "info"
	WarnLevelString  = "warn"
	ErrorLevelString = "error"
	FatalLevelString = "fatal"
)

// DefaultTimeFormat of logger
const DefaultTimeFormat = time.RFC3339

// StringToLevel to set string to level
func StringToLevel(level string) Level {
	switch strings.ToLower(level) {
	case DebugLevelString:
		return DebugLevel
	case InfoLevelString:
		return InfoLevel
	case WarnLevelString:
		return WarnLevel
	case ErrorLevelString:
		return ErrorLevel
	case FatalLevelString:
		return FatalLevel
	default:
		// TODO: make this more informative when happened
		return InfoLevel
	}
}

// LevelToString convert log level to readable string
func LevelToString(l Level) string {
	switch l {
	case DebugLevel:
		return DebugLevelString
	case InfoLevel:
		return InfoLevelString
	case WarnLevel:
		return WarnLevelString
	case ErrorLevel:
		return ErrorLevelString
	case FatalLevel:
		return FatalLevelString
	default:
		return InfoLevelString
	}
}

// CreateLogFile create a file and return io.Writer for file manipulation
func CreateLogFile(filename string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(filename), 0744)
	if err != nil && err != os.ErrExist {
		return nil, err
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}
