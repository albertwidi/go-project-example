package logrus

import (
	"io"
	"os"
	"sync"

	"github.com/albertwidi/go-project-example/internal/pkg/log/logger"
	"github.com/sirupsen/logrus"
)

var _ logger.Logger = (*Logger)(nil)

// Logger struct
type Logger struct {
	logger *logrus.Logger
	config logger.Config
	mu     *sync.Mutex
}

// DefaultLogger return default value of logger
func DefaultLogger() *Logger {
	l := Logger{
		config: logger.Config{
			Level:      logger.InfoLevel,
			TimeFormat: logger.DefaultTimeFormat,
		},
	}

	lgr := logrus.New()
	lgr.SetFormatter(&logrus.TextFormatter{
		DisableColors:   !l.config.UseColor,
		TimestampFormat: l.config.TimeFormat,
	})
	setLevel(lgr, l.config.Level)

	l.logger = lgr
	return &l
}

// New logger
func New(config *logger.Config) (*Logger, error) {
	if config == nil {
		config = &logger.Config{
			Level:      logger.InfoLevel,
			TimeFormat: logger.DefaultTimeFormat,
		}
	}

	if config.TimeFormat == "" {
		config.TimeFormat = logger.DefaultTimeFormat
	}

	lgr, err := newLogger(config)
	if err != nil {
		return nil, err
	}
	l := Logger{
		logger: lgr,
		config: *config,
	}
	return &l, nil
}

func newLogger(config *logger.Config) (*logrus.Logger, error) {
	lgr := logrus.New()

	// set writer to file if config.LogFile is not empty
	if config.LogFile != "" {
		file, err := logger.CreateLogFile(config.LogFile)
		if err != nil {
			return nil, err
		}
		// use multiwriter to write to os.Stderr and file output
		writers := io.MultiWriter(os.Stderr, file)
		lgr.SetOutput(writers)
	}

	// choose between json formatter or normal text
	if config.UseJSON {
		lgr.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: config.TimeFormat,
		})
	} else {
		lgr.SetFormatter(&logrus.TextFormatter{
			// invert the bool
			DisableColors:   !config.UseColor,
			TimestampFormat: config.TimeFormat,
		})
	}

	// set caller
	lgr.SetReportCaller(config.Caller)
	// set level
	setLevel(lgr, config.Level)

	return lgr, nil
}

// SetConfig to paply a new config to logger
func (l *Logger) SetConfig(config *logger.Config) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if config == nil {
		return nil
	}

	logger, err := newLogger(config)
	if err != nil {
		return err
	}

	l.logger = logger
	return nil
}

func setLevel(lgr *logrus.Logger, level logger.Level) {
	switch level {
	case logger.DebugLevel:
		lgr.SetLevel(logrus.DebugLevel)
	case logger.InfoLevel:
		lgr.SetLevel(logrus.InfoLevel)
	case logger.WarnLevel:
		lgr.SetLevel(logrus.WarnLevel)
	case logger.ErrorLevel:
		lgr.SetLevel(logrus.ErrorLevel)
	case logger.FatalLevel:
		lgr.SetLevel(logrus.FatalLevel)
	default:
		lgr.SetLevel(logrus.InfoLevel)
	}
}

// SetLevel for setting log level
func (l *Logger) SetLevel(level logger.Level) error {
	if level < logger.DebugLevel || level > logger.FatalLevel {
		level = logger.InfoLevel
	}

	if level != l.config.Level {
		setLevel(l.logger, level)
		l.config.Level = level
	}

	return nil
}

// Debug function
func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Debugf function
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
}

// Debugln function
func (l *Logger) Debugln(args ...interface{}) {
	l.logger.Debugln(args...)
}

// Debugw function
func (l *Logger) Debugw(message string, fields logger.KV) {
	l.logger.WithFields(logrus.Fields(fields)).Debugln(message)
}

// Info function
func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof function
func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}

// Infoln function
func (l *Logger) Infoln(args ...interface{}) {
	l.logger.Infoln(args...)
}

// Infow function
func (l *Logger) Infow(message string, fields logger.KV) {
	l.logger.WithFields(logrus.Fields(fields)).Infoln(message)
}

// Warn function
func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Warnf function
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger.Warnf(format, v...)
}

// Warnln function
func (l *Logger) Warnln(args ...interface{}) {
	l.logger.Warnln(args...)
}

// Warnw function
func (l *Logger) Warnw(message string, fields logger.KV) {
	l.logger.WithFields(logrus.Fields(fields)).Warnln(message)
}

// Error function
func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Errorf function
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Errorf(format, v...)
}

// Errorln function
func (l *Logger) Errorln(args ...interface{}) {
	l.logger.Errorln(args...)
}

// Errorw function
func (l *Logger) Errorw(message string, fields logger.KV) {
	l.logger.WithFields(logrus.Fields(fields)).Errorln(message)
}

// Fatal function
func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalf function
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatalf(format, v...)
}

// Fatalln function
func (l *Logger) Fatalln(args ...interface{}) {
	l.logger.Fatalln(args...)
}

// Fatalw function
func (l *Logger) Fatalw(message string, fields logger.KV) {
	l.logger.WithFields(logrus.Fields(fields)).Fatalln(message)
}
