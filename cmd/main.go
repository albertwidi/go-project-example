package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/albertwidi/kothak/cmd/project"
)

func main() {
	f := project.Flags{}
	flag.StringVar(&f.ConfigurationFile, "config_file", "./project_config.yaml", "configuration file of the project")
	flag.Var(&f.EnvironmentFiles, "env_files", "environment files for project configuration")
	flag.StringVar(&f.TimeZone, "tz", "Asia/Jakarta", "time zone of the project")
	flag.BoolVar(&f.DebugMode, "debug", false, "turn on debug mode, this will set log level to debug")
	flag.Parse()

	if err := project.Run(f); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}
