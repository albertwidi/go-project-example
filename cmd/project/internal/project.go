package project

import (
	"github.com/albertwidi/go_project_example/internal/config"
	"github.com/albertwidi/go_project_example/internal/pkg/log/logger"
)

type arrayFlags []string

// String return string implementation of array flags
func (af arrayFlags) String() string {
	return ""
}

// Set for append the value of arrayFlags
func (af arrayFlags) Set(value string) error {
	af = append(af, value)
	return nil
}

// Flags of project
type Flags struct {
	DebugMode         bool
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
	return nil
}
