package cucumber

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

// Feature interface
type Feature interface {
	BeforeRegister() error
	SetLogger(logger *log.Logger)
	FeatureContext(s *godog.Suite)
}

// Cucumber object
type Cucumber struct {
	features []Feature
	options  Options
	// logger
	logger *log.Logger
}

// Options struct
type Options struct {
	Debug Debug
}

// Debug options
type Debug struct {
	LogFile string
}

// New cucumber instance
func New(opts *Options) (*Cucumber, error) {
	var (
		options Options
		f       *os.File
		err     error
	)
	if opts != nil {
		options = *opts
	}

	// set the logger
	if options.Debug.LogFile != "" {
		_, err := os.Stat(options.Debug.LogFile)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}

		if os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Dir(options.Debug.LogFile), 0744)
			if err != nil && err != os.ErrExist {
				return nil, err
			}

			f, err = os.OpenFile(options.Debug.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				return nil, err
			}
		}
		if err == nil {
			f, err = os.OpenFile(options.Debug.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				return nil, err
			}
		}
	} else {
		// throw logger to dev null
		f, err = os.Open("/dev/null")
		if err != nil {
			return nil, err
		}
	}

	// write to multi writer
	multi := io.MultiWriter(os.Stdout, f)
	// set logger
	logger := &log.Logger{}
	logger.SetOutput(multi)
	c := Cucumber{
		features: make([]Feature, 0),
		options:  options,
		logger:   logger,
	}
	return &c, nil
}

// Logger return the cucumber logger
func (c *Cucumber) Logger() *log.Logger {
	return c.logger
}

// RegisterFeatures func
func (c *Cucumber) RegisterFeatures(features ...Feature) error {
	for _, f := range features {
		f.SetLogger(c.Logger())
		if err := f.BeforeRegister(); err != nil {
			return err
		}
	}
	c.features = append(c.features, features...)
	return nil
}

// FeatureContext for triggering godog
func (c *Cucumber) FeatureContext(s *godog.Suite) {
	for _, f := range c.features {
		f.FeatureContext(s)
	}
}
