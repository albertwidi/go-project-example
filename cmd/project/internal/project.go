package project

import (
	"context"

	"github.com/albertwidi/go_project_example/internal/config"
	"github.com/albertwidi/go_project_example/internal/pkg/kothak"
	"github.com/albertwidi/go_project_example/internal/pkg/log/logger"
)

type arrayFlags []string

// String return string implementation of array flags
func (af *arrayFlags) String() string {
	return ""
}

// Set for append the value of arrayFlags
func (af *arrayFlags) Set(value string) error {
	*af = append(*af, value)
	return nil
}

// Flags of project
type Flags struct {
	Debug             debugFlag
	Log               logFlag
	Dev               bool
	TimeZone          string
	ConfigurationFile string
	EnvironmentFiles  arrayFlags
	LogFile           string
}

// Config of project
type Config struct {
	config.DefaultConfig
}

// Run the project
func Run(f Flags, logger logger.Logger) error {
	// load project configuration
	projectConfig := Config{}
	if err := config.ParseFile(f.ConfigurationFile, &projectConfig, f.EnvironmentFiles...); err != nil {
		return err
	}

	if f.Debug.TestConfig {
		logger.Infof("testing config with flags and configurations:")
		logger.Infof("flags:\n%+v", f)
		logger.Infof("config:\n%+v", projectConfig)
	}

	resources, err := kothak.New(context.TODO(), projectConfig.Resources, logger)
	if err != nil {
		return err
	}
	// close all connections when program exiting
	defer resources.CloseAll()

	// exit early if we only test config, do not run the server
	if f.Debug.TestConfig {
		return nil
	}

	return nil
}
