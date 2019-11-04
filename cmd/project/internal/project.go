package project

import (
	"context"

	"github.com/albertwidi/go_project_example/internal/config"
	"github.com/albertwidi/go_project_example/internal/pkg/kothak"
	lg "github.com/albertwidi/go_project_example/internal/pkg/log/logger"
	"github.com/albertwidi/go_project_example/internal/pkg/log/logger/zap"
)

// Flags of project
type Flags struct {
	Debug             debugFlag
	Log               logFlag
	EnvironmentFiles  envFileFlag
	TimeZone          string
	ConfigurationFile string
	LogFile           string
}

// Config of project
type Config struct {
	config.DefaultConfig
}

// Run the project
func Run(f Flags) error {
	// initiate project logger
	logger, err := zap.New(&lg.Config{
		Level:    lg.StringToLevel(f.Log.Level),
		LogFile:  f.Log.File,
		UseColor: f.Log.Color,
	})
	if err != nil {
		return err
	}

	// load project configuration
	projectConfig := Config{}
	if err := config.ParseFile(f.ConfigurationFile, &projectConfig, f.EnvironmentFiles.envFiles...); err != nil {
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
