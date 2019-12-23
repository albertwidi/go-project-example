package main

import (
	"flag"
	"fmt"
	"os"

	project "github.com/albertwidi/go-project-example/cmd/project/internal"
)

func main() {
	exitCode := 0
	f := project.Flags{}
	flag.StringVar(&f.ConfigurationFile, "config_file", "./aha.config.toml", "configuration file of the project")
	flag.Var(&f.EnvironmentFile, "env_file", "helper file for environment variable configuration")
	flag.StringVar(&f.TimeZone, "tz", "", "time zone of the project")
	flag.Var(&f.Debug, "debug", "turn on debug mode, this will set log level to debug")
	flag.Parse()

	if err := project.Run(f); err != nil {
		exitCode = 1
		fmt.Printf("%v", err)
	}
	os.Exit(exitCode)
}
