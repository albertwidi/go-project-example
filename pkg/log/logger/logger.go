// logging library using uber zap

package logger

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/albertwidi/kothak/pkg/errors"
	"go.uber.org/zap"
)

// Level type
type Level int

// level of log
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

// Config of logger
type Config struct {
	Level      Level
	Fileoutput string
}

// Log struct
type Log struct {
	logger  *zap.Logger
	sugared *zap.SugaredLogger
	zapConf *zap.Config
	// log configuration
	config Config
}

func newZapConfig() *zap.Config {
	config := zap.NewProductionConfig()
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.ErrorOutputPaths = []string{"stderr"}
	return &config
}

func stringToLevel(s string) Level {
	switch strings.ToLower(s) {
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

func levelToString(l Level) string {
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

// New logger
func New(conf Config) (*Log, error) {
	if conf.Level == Level(0) {
		conf.Level = InfoLevel
	}

	l := Log{
		zapConf: newZapConfig(),
		config:  conf,
	}
	l.setLevel()
	l.setOutputToFile()
	// build all the configuration after level and file output
	if err := l.build(); err != nil {
		return nil, err
	}
	return &l, nil
}

// Default logger
func Default() *Log {
	l := Log{
		zapConf: newZapConfig(),
		config: Config{
			Level: InfoLevel,
		},
	}
	l.setLevel()
	l.build()
	return &l
}

func (l *Log) build() error {
	// attempted to create log file, if not created then the logger will return error
	if l.config.Fileoutput != "" {
		path, _ := filepath.Split(l.config.Fileoutput)
		if path != "" {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				err := os.MkdirAll(path, 0755)
				if err != nil {
					return err
				}
			}
		}
		// check if the log file is exists
		if _, err := os.Stat(l.config.Fileoutput); os.IsNotExist(err) {
			if _, err = os.Create(l.config.Fileoutput); err != nil {
				log.Println("ERROR HERE")
				return err
			}
		}
	}

	zapLogger, err := l.zapConf.Build()
	if err != nil {
		return err
	}
	l.logger = zapLogger
	l.logger.Sync()
	l.sugared = l.logger.Sugar()
	l.sugared.Sync()
	return nil
}

func (l *Log) setLevel() {
	switch l.config.Level {
	case DebugLevel:
		l.zapConf.Level.SetLevel(zap.DebugLevel)
	case InfoLevel:
		l.zapConf.Level.SetLevel(zap.InfoLevel)
	case WarnLevel:
		l.zapConf.Level.SetLevel(zap.WarnLevel)
	case ErrorLevel:
		l.zapConf.Level.SetLevel(zap.ErrorLevel)
	case FatalLevel:
		l.zapConf.Level.SetLevel(zap.FatalLevel)
	default:
		l.zapConf.Level.SetLevel(zap.InfoLevel)
	}
}

func (l *Log) setOutputToFile() {
	if l.config.Fileoutput == "" {
		return
	}
	l.zapConf.OutputPaths = []string{"stdout", l.config.Fileoutput}
}

// SetLevel to logger
func (l *Log) SetLevel(lvl Level) {
	l.config.Level = lvl
	l.setLevel()
	l.build()
}

// SetLevelString set level using string instead of level
func (l *Log) SetLevelString(lvl string) {
	l.SetLevel(stringToLevel(lvl))
}

// GetLevel of logger
func (l *Log) GetLevel() string {
	return levelToString(l.config.Level)
}

// SetOutputToFile to logger
func (l *Log) SetOutputToFile(filepath string) error {
	l.config.Fileoutput = filepath
	l.setOutputToFile()
	return l.build()
}

// Debug log
func (l *Log) Debug(args ...interface{}) {
	l.sugared.Debug(args...)
}

// Debugf log
func (l *Log) Debugf(format string, args ...interface{}) {
	l.sugared.Debugf(format, args...)
}

// Debugw log
func (l *Log) Debugw(msg string, keyAndValues ...interface{}) {
	l.sugared.Debugw(msg, keyAndValues...)
}

// Print log
func (l *Log) Print(args ...interface{}) {
	l.sugared.Info(args...)
}

// Println log
func (l *Log) Println(args ...interface{}) {
	l.sugared.Info(args...)
}

// Printf log
func (l *Log) Printf(format string, args ...interface{}) {
	l.sugared.Infof(format, args...)
}

// Printw log
func (l *Log) Printw(msg string, keyAndValues ...interface{}) {
	l.sugared.Infow(msg, keyAndValues...)
}

// Info log
func (l *Log) Info(args ...interface{}) {
	l.sugared.Info(args...)
}

// Infof log
func (l *Log) Infof(format string, args ...interface{}) {
	l.sugared.Infof(format, args...)
}

// Infow log
func (l *Log) Infow(msg string, keyAndValues ...interface{}) {
	l.sugared.Infow(msg, keyAndValues...)
}

// Warn log
func (l *Log) Warn(args ...interface{}) {
	l.sugared.Warn(args...)
}

// Warnf log
func (l *Log) Warnf(format string, args ...interface{}) {
	l.sugared.Warnf(format, args...)
}

// Warnw log
func (l *Log) Warnw(msg string, keyAndValues ...interface{}) {
	l.sugared.Warnw(msg, keyAndValues...)
}

// Error log
func (l *Log) Error(args ...interface{}) {
	l.sugared.Error(args...)
}

// Errorf log
func (l *Log) Errorf(format string, args ...interface{}) {
	l.sugared.Errorf(format, args...)
}

// Errorw log
func (l *Log) Errorw(msg string, keyAndValues ...interface{}) {
	l.sugared.Errorw(msg, keyAndValues...)
}

// Errors log log error detail from Errs
func (l *Log) Errors(err error) {
	l.sugared.With(err).Error(err.Error())
}

// Fatal log
func (l *Log) Fatal(args ...interface{}) {
	l.sugared.Fatal(args...)
}

// Fatalf log
func (l *Log) Fatalf(format string, args ...interface{}) {
	l.sugared.Fatalf(format, args...)
}

// Fatalw log
func (l *Log) Fatalw(format string, keyAndValues ...interface{}) {
	l.sugared.Fatalw(format, keyAndValues...)
}

// With log
func (l *Log) With(args ...interface{}) *zap.SugaredLogger {
	var v []interface{}
	for _, a := range args {
		switch a.(type) {
		case *errors.Error:
			v = append(v, errorsFieldsIntf(a.((*errors.Error)))...)
		default:
			v = append(v, a)
		}
	}
	return l.sugared.With(args...)
}

func errorsFieldsIntf(e *errors.Error) []interface{} {
	return e.GetFields().ToArrayInterface()
}

// Config return log configuration
func (l *Log) Config() Config {
	return l.config
}
