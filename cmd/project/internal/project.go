package project

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/albertwidi/go-project-example/internal/config"
	"github.com/albertwidi/go-project-example/internal/kothak"
	lg "github.com/albertwidi/go-project-example/internal/pkg/log/logger"
	"github.com/albertwidi/go-project-example/internal/pkg/log/logger/zap"
	"github.com/albertwidi/go-project-example/internal/server"
)

// Flags of project
type Flags struct {
	Debug             debugFlag
	EnvironmentFile   envFileFlag
	TimeZone          string
	ConfigurationFile string
	LogFile           string
	Version           bool
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
		return fmt.Errorf("run: error when initiating logger: %w", err)
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

	s, err := server.New(projectConfig.Servers.Admin.Address, debugServer)
	if err != nil {
		return err
	}
	// run the server
	errChan := s.Run()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	// exit early if we only test config
	testChan := make(chan struct{}, 1)
	if f.Debug.TestConfig {
		logger.Infoln("testing: giving time for server to run")
		go func() {
			time.Sleep(time.Second * 5)
			testChan <- struct{}{}
		}()
	}

	select {
	case err := <-errChan:
		return err
	case sig := <-sigChan:
		switch sig {
		case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
			return errors.New("project: receive signal to terminate program")
		}
	case <-testChan:
		logger.Infoln("testing: test completed successfully")
		return nil
	}
	return nil
}

func newMainServer() {

}
