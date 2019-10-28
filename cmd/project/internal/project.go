package project

import (
	"fmt"
	"time"

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

// Run the project
func Run(f Flags, logger logger.Logger) error {
	fmt.Println(time.Now())
	return nil
}
