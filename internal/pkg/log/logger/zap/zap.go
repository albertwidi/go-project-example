package zap

import (
	"github.com/albertwidi/go-project-example/internal/pkg/log/logger"
	"go.uber.org/zap"
	"sync"
)

var _ logger.Logger = (*Logger)(nil)

// Logger struct
type Logger struct {
	logger    *zap.Logger
	sugared   *zap.SugaredLogger
	zapconfig zap.Config
	config    *logger.Config
	mu        sync.Mutex
}

// New zap logger
func New(config *logger.Config) (*Logger, error) {
	if config == nil {
		config = &logger.Config{
			Level:      logger.InfoLevel,
			TimeFormat: logger.DefaultTimeFormat,
		}
	}

	l := Logger{}
	if err := l.initLogger(config); err != nil {
		return nil, err
	}
	return &l, nil
}

// initLogger is mimicking zap.NewProductionConfig()
// to standarize our own internal configuration, the function means to provide the production ready configuration
func (l *Logger) initLogger(config *logger.Config) error {
	encoding := "console"
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	if config.UseJSON {
		encoding = "json"
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	errAndOutputPaths := []string{"stderr"}
	if config.LogFile != "" {
		if _, err := logger.CreateLogFile(config.LogFile); err != nil {
			return err
		}
		errAndOutputPaths = append(errAndOutputPaths, config.LogFile)
	}

	zapConfig := zap.Config{
		Level:         getLevel(config.Level),
		Encoding:      encoding,
		EncoderConfig: encoderConfig,
		// set disable to true when config.Caller is false
		DisableCaller:    config.Caller == false,
		OutputPaths:      errAndOutputPaths,
		ErrorOutputPaths: errAndOutputPaths,
	}

	z, err := zapConfig.Build()
	if err != nil {
		return err
	}
	l.logger = z
	l.sugared = z.Sugar()
	l.zapconfig = zapConfig
	l.config = config
	return nil
}

func getLevel(level logger.Level) zap.AtomicLevel {
	var atl zap.AtomicLevel

	switch level {
	case logger.DebugLevel:
		atl = zap.NewAtomicLevelAt(zap.DebugLevel)
	case logger.InfoLevel:
		atl = zap.NewAtomicLevelAt(zap.InfoLevel)
	case logger.WarnLevel:
		atl = zap.NewAtomicLevelAt(zap.WarnLevel)
	case logger.ErrorLevel:
		atl = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case logger.FatalLevel:
		atl = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		atl = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	return atl
}

// SetLevel for setting log level
func (l *Logger) SetLevel(level logger.Level) error {
	if level < logger.DebugLevel || level > logger.FatalLevel {
		level = logger.InfoLevel
	}

	if level != l.config.Level {
		l.config.Level = level
		if err := l.initLogger(l.config); err != nil {
			return err
		}
	}
	return nil
}

// SetConfig to paply a new config to logger
func (l *Logger) SetConfig(config *logger.Config) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if config == nil {
		return nil
	}

	if err := l.initLogger(config); err != nil {
		return err
	}
	return nil
}

// Debug function
func (l *Logger) Debug(args ...interface{}) {
	l.sugared.Debug(args...)
}

// Debugf function
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.sugared.Debugf(format, v...)
}

// Debugln function
func (l *Logger) Debugln(args ...interface{}) {
	l.sugared.Debug(args...)
}

// Debugw function
func (l *Logger) Debugw(message string, fields logger.KV) {
	l.sugared.Debugw(message, fieldsToKV(fields))
}

// Info function
func (l *Logger) Info(args ...interface{}) {
	l.sugared.Info(args...)
}

// Infof function
func (l *Logger) Infof(format string, v ...interface{}) {
	l.sugared.Infof(format, v...)
}

// Infoln function
func (l *Logger) Infoln(args ...interface{}) {
	l.sugared.Info(args...)
}

// Infow function
func (l *Logger) Infow(message string, fields logger.KV) {
	l.sugared.Infow(message, fieldsToKV(fields))
}

// Warn function
func (l *Logger) Warn(args ...interface{}) {
	l.sugared.Warn(args...)
}

// Warnf function
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.sugared.Warnf(format, v...)
}

// Warnln function
func (l *Logger) Warnln(args ...interface{}) {
	l.sugared.Warn(args...)
}

// Warnw function
func (l *Logger) Warnw(message string, fields logger.KV) {
	l.sugared.Warnw(message, fieldsToKV(fields))
}

// Error function
func (l *Logger) Error(args ...interface{}) {
	l.sugared.Error(args...)
}

// Errorf function
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.sugared.Errorf(format, v...)
}

// Errorln function
func (l *Logger) Errorln(args ...interface{}) {
	l.sugared.Error(args...)
}

// Errorw function
func (l *Logger) Errorw(message string, fields logger.KV) {
	l.sugared.Errorw(message, fieldsToKV(fields))
}

// Fatal function
func (l *Logger) Fatal(args ...interface{}) {
	l.sugared.Fatal(args...)
}

// Fatalf function
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.sugared.Fatalf(format, v...)
}

// Fatalln function
func (l *Logger) Fatalln(args ...interface{}) {
	l.sugared.Fatal(args...)
}

// Fatalw function
func (l *Logger) Fatalw(message string, fields logger.KV) {
	l.sugared.Fatalw(message, fieldsToKV(fields))
}

func fieldsToKV(fields logger.KV) []interface{} {
	var i int
	kv := make([]interface{}, len(fields)*2)
	for k, v := range fields {
		kv[i] = k
		kv[i+1] = v
		i += 2
	}
	return kv
}
