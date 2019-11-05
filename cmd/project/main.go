package main

import (
	"flag"
	"fmt"
	"os"

	project "github.com/albertwidi/go_project_example/cmd/project/internal"
)

func main() {
	f := project.Flags{}
	flag.StringVar(&f.ConfigurationFile, "config_file", "./aha.config.toml", "configuration file of the project")
	flag.Var(&f.EnvironmentFile, "env_file", "helper file for environment variable configuration")
	flag.Var(&f.Log, "log", "log configuration")
	flag.StringVar(&f.TimeZone, "tz", "Asia/Jakarta", "time zone of the project")
	flag.Var(&f.Debug, "debug", "turn on debug mode, this will set log level to debug")
	flag.Parse()

	if err := project.Run(f); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}
