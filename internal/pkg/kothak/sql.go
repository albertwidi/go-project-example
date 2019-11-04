package kothak

import (
	"time"

	"github.com/albertwidi/go_project_example/internal/pkg/defaults"
)

// DBConfig define sql databases configuration
type DBConfig struct {
	MaxRetry              int           `yaml:"max_retry" toml:"max_retry" default:"1"`
	MaxOpenConnections    int           `yaml:"max_open_conns" toml:"max_open_conns" default:"10"`
	MaxIdleConnections    int           `yaml:"max_idle_conns" toml:"max_idle_conns" default:"2"`
	ConnectionMaxLifetime string        `yaml:"conn_max_lifetime" toml:"conn_max_lifetime" default:"30s"`
	SQLDBs                []SQLDBConfig `yaml:"connect" toml:"connect"`
	connMaxLifetime       time.Duration
}

// SetDefault configuration
func (dbconf *DBConfig) SetDefault() error {
	if err := defaults.SetDefault(dbconf); err != nil {
		return err
	}

	if dbconf.ConnectionMaxLifetime != "" {
		dur, err := time.ParseDuration(dbconf.ConnectionMaxLifetime)
		if err != nil {
			return err
		}
		dbconf.connMaxLifetime = dur
	}
	return nil
}

// SQLDBConfig of kothak
type SQLDBConfig struct {
	Name              string                `yaml:"name" toml:"name"`
	Driver            string                `yaml:"driver" toml:"driver"`
	LeaderConnConfig  SQLDBConnectionConfig `yaml:"leader" toml:"leader"`
	ReplicaConnConfig SQLDBConnectionConfig `yaml:"replica" toml:"replica"`
}

// SQLDBConnectionConfig struct
type SQLDBConnectionConfig struct {
	DSN                   string `yaml:"dsn" toml:"dsn"`
	MaxOpenConnections    int    `yaml:"max_open_conns" toml:"max_open_conns"`
	MaxIdleConnections    int    `yaml:"max_idle_conns" toml:"max_idle_conns"`
	ConnectionMaxLifetime string `yaml:"conn_max_lifetime" toml:"conn_max_lifetime"`
	MaxRetry              int    `yaml:"max_retry" toml:"max_retry"`
	connMaxLifeTime       time.Duration
}

// SetDefault configuration
func (connConfig *SQLDBConnectionConfig) SetDefault(dbconfig DBConfig) error {
	if connConfig.MaxRetry == 0 {
		connConfig.MaxRetry = dbconfig.MaxRetry
	}

	if connConfig.MaxOpenConnections == 0 {
		connConfig.MaxOpenConnections = dbconfig.MaxOpenConnections
	}

	if connConfig.MaxIdleConnections == 0 {
		connConfig.MaxIdleConnections = dbconfig.MaxIdleConnections
	}

	if connConfig.ConnectionMaxLifetime == "" {
		connConfig.ConnectionMaxLifetime = dbconfig.ConnectionMaxLifetime
	}

	connConfig.connMaxLifeTime = dbconfig.connMaxLifetime
	if connConfig.ConnectionMaxLifetime != "" {
		dur, err := time.ParseDuration(connConfig.ConnectionMaxLifetime)
		if err != nil {
			return err
		}
		connConfig.connMaxLifeTime = dur
	}
	return nil
}
