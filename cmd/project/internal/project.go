package project

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/albertwidi/go_project_example/internal/config"
	"github.com/albertwidi/go_project_example/internal/kothak"
	lg "github.com/albertwidi/go_project_example/internal/pkg/log/logger"
	"github.com/albertwidi/go_project_example/internal/pkg/log/logger/zap"
	"github.com/albertwidi/go_project_example/internal/server"
)

// Flags of project
type Flags struct {
	Debug             debugFlag
	EnvironmentFile   envFileFlag
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
	// set default timezone
	os.Setenv("TZ", f.TimeZone)

	// load project configuration
	projectConfig := Config{}
	if err := config.ParseFile(f.ConfigurationFile, &projectConfig, f.EnvironmentFile.envFiles...); err != nil {
		return err
	}

	// initiate project logger
	logger, err := zap.New(&lg.Config{
		Level:    lg.StringToLevel(projectConfig.Log.Level),
		LogFile:  projectConfig.Log.File,
		UseColor: projectConfig.Log.Color,
	})
	if err != nil {
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

	// repositories
	repo, err := newRepositories(resources)
	if err != nil {
		return err
	}

	// initiate new servers
	newMainServer()
	debugServerAddr := projectConfig.Servers.Debug.Address
	debugServer, err := newDebugServer(debugServerAddr, repo)
	if err != nil {
		return err
	}
	newAdminServer()

	s, err := server.New(debugServer)
	if err != nil {
		return err
	}

	errChan := make(chan error, 1)
	errChan <- s.Run()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case err := <-errChan:
		return err
	case sig := <-sigChan:
		switch sig {
		case syscall.SIGTERM, syscall.SIGQUIT:
			return errors.New("project: receiving signal to terminate program")
		}
	}
	return nil
}

func newAdminServer() {

}

func newMainServer() {

}
