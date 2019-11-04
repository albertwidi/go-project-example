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
	logger.Infof("%+v", f)
	// load project configuration
	projectConfig := Config{}
	if err := config.ParseFile(f.ConfigurationFile, &projectConfig, f.EnvironmentFiles...); err != nil {
		return err
	}

	logger.Infof("%+v", projectConfig)

	resources, err := kothak.New(context.TODO(), projectConfig.Resources, logger)
	if err != nil {
		return err
	}

	resources.CloseAll()
	return nil
}
