// using go standard logger

package std

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/albertwidi/go-project-example/internal/pkg/log/logger"
)

var _ logger.Logger = (*Logger)(nil)

// Logger of go standard logger
type Logger struct {
	logger *log.Logger
	config *logger.Config
}

var levelFormat = []string{
	"[DEBUG]",
	"[INFO]",
	"[WARN]",
	"[ERROR]",
	"[FATAL]",
}

// New standard logger
func New(config *logger.Config) (*Logger, error) {
	return nil, nil
}

func newLogger(config *logger.Config) (*log.Logger, error) {
	stdLogger := &log.Logger{}

	if config == nil {
		config = &logger.Config{
			Level:      logger.InfoLevel,
			TimeFormat: logger.DefaultTimeFormat,
		}
		stdLogger.SetOutput(os.Stderr)
	}

	if config.TimeFormat == "" {
		config.TimeFormat = logger.DefaultTimeFormat
	}

	if config.LogFile != "" {
		file, err := logger.CreateLogFile(config.LogFile)
		if err != nil {
			return nil, err
		}
		writers := io.MultiWriter(os.Stderr, file)
		stdLogger.SetOutput(writers)
	}

	return stdLogger, nil
}

// SetConfig to reset logger configuration
func (l *Logger) SetConfig(config *logger.Config) error {
	return nil
}

// SetLevel to set log level
func (l *Logger) SetLevel(level logger.Level) error {
	return nil
}

// Debug log using standard logger
func (l *Logger) Debug(args ...interface{}) {
	l.print(logger.DebugLevel, args...)
}

// Debugf log using standard logger
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.printf(logger.DebugLevel, format, v...)
}

// Debugln log using standard logger
func (l *Logger) Debugln(args ...interface{}) {
	l.println(logger.DebugLevel, args...)
}

// Debugw log using standard logger
func (l *Logger) Debugw(message string, fields logger.KV) {
	l.printw(logger.DebugLevel, message, fields)
}

// Info log using standard logger
func (l *Logger) Info(args ...interface{}) {
	l.print(logger.InfoLevel, args...)
}

// Infof log using standard logger
func (l *Logger) Infof(format string, v ...interface{}) {
	l.printf(logger.InfoLevel, format, v...)
}

// Infoln log using standard logger
func (l *Logger) Infoln(args ...interface{}) {
	l.println(logger.InfoLevel, args...)
}

// Infow log using standard logger
func (l *Logger) Infow(message string, fields logger.KV) {
	l.printw(logger.InfoLevel, message, fields)
}

// Warn log using standard logger
func (l *Logger) Warn(args ...interface{}) {
	l.print(logger.WarnLevel, args...)
}

// Warnf log using standard logger
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.printf(logger.WarnLevel, format, v...)
}

// Warnln log using standard logger
func (l *Logger) Warnln(args ...interface{}) {
	l.println(logger.WarnLevel, args...)
}

// Warnw log using standard logger
func (l *Logger) Warnw(message string, fields logger.KV) {
	l.printw(logger.WarnLevel, message, fields)
}

// Error log using standard logger
func (l *Logger) Error(args ...interface{}) {
	l.print(logger.ErrorLevel, args...)
}

// Errorf log using standard logger
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.printf(logger.ErrorLevel, format, v...)
}

// Errorln log using standard logger
func (l *Logger) Errorln(args ...interface{}) {
	l.println(logger.ErrorLevel, args...)
}

// Errorw log using standard logger
func (l *Logger) Errorw(message string, fields logger.KV) {
	l.printw(logger.ErrorLevel, message, fields)
}

// Fatal log using standard logger
func (l *Logger) Fatal(args ...interface{}) {
	l.print(logger.FatalLevel, args...)
	os.Exit(1)
}

// Fatalf log using standard logger
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.printf(logger.FatalLevel, format, v...)
	os.Exit(1)
}

// Fatalln log using standard logger
func (l *Logger) Fatalln(args ...interface{}) {
	l.println(logger.FatalLevel, args...)
	os.Exit(1)
}

// Fatalw log using standard logger
func (l *Logger) Fatalw(message string, fields logger.KV) {
	l.printw(logger.FatalLevel, message, fields)
	os.Exit(1)
}

func (l *Logger) print(level logger.Level, args ...interface{}) {
	if level < l.config.Level {
		return
	}
	l.logger.Print(levelFormat[level], args)
}

func (l *Logger) printf(level logger.Level, format string, v ...interface{}) {
	if level < l.config.Level {
		return
	}
	format = levelFormat[level] + format
	l.logger.Printf(format, v)
}

func (l *Logger) println(level logger.Level, args ...interface{}) {
	if level < l.config.Level {
		return
	}
	l.logger.Println(levelFormat[level], args)
}

func (l *Logger) printw(level logger.Level, message string, fields logger.KV) {
	if level < l.config.Level {
		return
	}
	l.logger.Println(levelFormat[level], message, formatFields(fields))
}

func formatFields(fields logger.KV) string {
	fieldsStr := ""
	fieldsLen := len(fields)
	counter := 0

	for k, v := range fields {
		fieldsStr += fmt.Sprintf("%s=%v", k, v)
		counter++
		if counter < fieldsLen {
			fieldsStr += " "
		}
	}

	return fieldsStr
}
