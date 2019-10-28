package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/albertwidi/kothak/cmd/project/internal"
	"github.com/albertwidi/kothak/lib/log/logger"
	"github.com/albertwidi/kothak/lib/log/logger/zap"
)

func main() {
	f := internal.Flags{}
	flag.StringVar(&f.ConfigurationFile, "config_file", "./project_config.yaml", "configuration file of the project")
	flag.Var(&f.EnvironmentFiles, "env_files", "environment files for project configuration")
	flag.StringVar(&f.LogFile, "log_file", "", "log file output")
	flag.StringVar(&f.TimeZone, "tz", "Asia/Jakarta", "time zone of the project")
	flag.BoolVar(&f.DebugMode, "debug", false, "turn on debug mode, this will set log level to debug")
	flag.Parse()

	// set default timezone
	os.Setenv("TZ", f.TimeZone)

	logger, err := zap.New(&logger.Config{
		Level:    logger.InfoLevel,
		LogFile:  f.LogFile,
		UseColor: true,
	})
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	if err := internal.Run(f, logger); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}
