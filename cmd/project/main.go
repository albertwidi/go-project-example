package main

import (
	"flag"
	"fmt"
	"os"

	project "github.com/albertwidi/go-project-example/cmd/project/internal"
)

var (
	buildVersion string
	buildCommit  string
)

const (
	usage = `Usage:
	backend -config_file=./project.config.toml \
		-env_file=./project.env.toml
	`
)

func main() {
	exitCode := 0
	f := project.Flags{}
	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }
	flag.StringVar(&f.ConfigurationFile, "config_file", "./aha.config.toml", "configuration file of the project")
	flag.Var(&f.EnvironmentFile, "env_file", "helper file for environment variable configuration")
	flag.StringVar(&f.TimeZone, "tz", "", "time zone of the project")
	flag.BoolVar(&f.Version, "version", false, "to print version of the prgoram")
	flag.Var(&f.Debug, "debug", "turn on debug mode, this will set log level to debug")
	flag.Parse()

	if f.Version {
		fmt.Fprintf(os.Stderr, "version: %s\ncommit: %s\n", buildVersion, buildCommit)
		return
	}
	if err := project.Run(f); err != nil {
		exitCode = 1
		fmt.Printf("%v", err)
	}
	os.Exit(exitCode)
}
