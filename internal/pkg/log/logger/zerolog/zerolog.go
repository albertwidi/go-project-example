package zerolog

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/albertwidi/go-project-example/internal/pkg/log/logger"
	"github.com/rs/zerolog"
)

var _ logger.Logger = (*Logger)(nil)

// Logger struct
type Logger struct {
	logger zerolog.Logger
	config logger.Config
	mu     sync.Mutex
}

// DefaultLogger return default value of logger
func DefaultLogger() *Logger {
	l := Logger{
		config: logger.Config{
			Level:      logger.InfoLevel,
			TimeFormat: logger.DefaultTimeFormat,
		},
	}

	lgr := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    !l.config.UseColor,
		TimeFormat: l.config.TimeFormat,
	})
	lgr = setLevel(lgr, l.config.Level)

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

func newLogger(config *logger.Config) (zerolog.Logger, error) {
	zerolog.TimeFieldFormat = config.TimeFormat
	zerolog.CallerSkipFrameCount = 4

	var writers zerolog.LevelWriter
	if config.UseJSON {
		writers = zerolog.MultiLevelWriter(os.Stderr)
	} else {
		writers = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			NoColor:    !config.UseColor,
			TimeFormat: config.TimeFormat,
		})
	}

	// set writer to file if config.LogFile is not empty
	if config.LogFile != "" {
		err := os.MkdirAll(filepath.Dir(config.LogFile), 0755)
		if err != nil {
			return zerolog.Logger{}, err
		}
		file, err := os.OpenFile(config.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return zerolog.Logger{}, err
		}
		writers = zerolog.MultiLevelWriter(writers, file)
	}

	lgr := zerolog.New(writers)
	lgr = setLevel(lgr, config.Level)
	if config.Caller {
		lgr = lgr.With().Caller().Logger()
	}
	return lgr, nil
}

// SetConfig to set a new logger configuration
func (l *Logger) SetConfig(config *logger.Config) error {
	l.mu.Lock()
	l.mu.Unlock()

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

func setLevel(lgr zerolog.Logger, level logger.Level) zerolog.Logger {
	switch level {
	case logger.DebugLevel:
		lgr = lgr.Level(zerolog.DebugLevel)
	case logger.InfoLevel:
		lgr = lgr.Level(zerolog.InfoLevel)
	case logger.WarnLevel:
		lgr = lgr.Level(zerolog.WarnLevel)
	case logger.ErrorLevel:
		lgr = lgr.Level(zerolog.ErrorLevel)
	case logger.FatalLevel:
		lgr = lgr.Level(zerolog.FatalLevel)
	default:
		lgr = lgr.Level(zerolog.InfoLevel)
	}
	return lgr
}

// SetLevel for setting log level
func (l *Logger) SetLevel(level logger.Level) error {
	if level < logger.DebugLevel || level > logger.FatalLevel {
		level = logger.InfoLevel
	}

	if level != l.config.Level {
		l.logger = setLevel(l.logger, level)
		l.config.Level = level
	}

	return nil
}

// Debug function
func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug().Timestamp().Msg(fmt.Sprint(args...))
}

// Debugf function
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger.Debug().Timestamp().Msgf(format, v...)
}

// Debugw function
func (l *Logger) Debugw(msg string, KV logger.KV) {
	l.logger.Debug().Timestamp().Fields(KV).Msg(msg)
}

// Info function
func (l *Logger) Info(args ...interface{}) {
	l.logger.Info().Timestamp().Msg(fmt.Sprint(args...))
}

// Infof function
func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Info().Timestamp().Msgf(format, v...)
}

// Infow function
func (l *Logger) Infow(msg string, KV logger.KV) {
	l.logger.Info().Timestamp().Fields(KV).Msg(msg)
}

// Warn function
func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn().Timestamp().Msg(fmt.Sprint(args...))
}

// Warnf function
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger.Warn().Timestamp().Msgf(format, v...)
}

// Warnw function
func (l *Logger) Warnw(msg string, KV logger.KV) {
	l.logger.Warn().Timestamp().Fields(KV).Msg(msg)
}

// Error function
func (l *Logger) Error(args ...interface{}) {
	l.logger.Error().Timestamp().Msg(fmt.Sprint(args...))
}

// Errorf function
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Error().Timestamp().Msgf(format, v...)
}

// Errorw function
func (l *Logger) Errorw(msg string, KV logger.KV) {
	l.logger.Error().Timestamp().Fields(KV).Msg(msg)
}

// Fatal function
func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal().Timestamp().Msg(fmt.Sprint(args...))
}

// Fatalf function
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal().Timestamp().Msgf(format, v...)
}

// Fatalw function
func (l *Logger) Fatalw(msg string, KV logger.KV) {
	l.logger.Fatal().Timestamp().Fields(KV).Msg(msg)
}
