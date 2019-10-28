package kothak

import (
	"time"
)

// DBConfig define sql databases configuration
type DBConfig struct {
	MaxRetry              int           `yaml:"max_retry"`
	MaxOpenConnections    int           `yaml:"max_open_conns"`
	MaxIdleConnections    int           `yaml:"max_idle_conns"`
	ConnectionMaxLifetime string        `yaml:"conn_max_lifetime"`
	SQLDBs                []SQLDBConfig `yaml:"connect"`

	connMaxLifetime time.Duration
}

// SetDefault configuration
func (dbconf *DBConfig) SetDefault() error {
	// check max retry default
	if dbconf.MaxRetry == 0 {
		dbconf.MaxRetry = 1
	}

	if dbconf.MaxOpenConnections == 0 {
		dbconf.MaxOpenConnections = 100
	}

	if dbconf.MaxIdleConnections == 0 {
		dbconf.MaxIdleConnections = 20
	}

	dbconf.connMaxLifetime = time.Second * 30
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
	Name                  string `yaml:"name"`
	Driver                string `yaml:"driver"`
	MasterDSN             string `yaml:"master"`
	FollowerDSN           string `yaml:"follower"`
	MaxOpenConnections    int    `yaml:"max_open_conns"`
	MaxIdleConnections    int    `yaml:"max_idle_conns"`
	ConnectionMaxLifetime string `yaml:"conn_max_lifetime"`
	MaxRetry              int    `yaml:"max_retry"`

	connMaxLifetime time.Duration
}

// SetDefault configuration
func (sqlconf *SQLDBConfig) SetDefault(dbconfig DBConfig) error {
	if sqlconf.MaxRetry == 0 {
		sqlconf.MaxRetry = dbconfig.MaxRetry
	}

	if sqlconf.MaxOpenConnections == 0 {
		sqlconf.MaxOpenConnections = dbconfig.MaxOpenConnections
	}

	if sqlconf.MaxIdleConnections == 0 {
		sqlconf.MaxIdleConnections = dbconfig.MaxIdleConnections
	}

	if sqlconf.ConnectionMaxLifetime == "" {
		sqlconf.ConnectionMaxLifetime = dbconfig.ConnectionMaxLifetime
	}

	if sqlconf.ConnectionMaxLifetime != "" {
		dur, err := time.ParseDuration(sqlconf.ConnectionMaxLifetime)
		if err != nil {
			return err
		}
		sqlconf.connMaxLifetime = dur
	} else {
		sqlconf.connMaxLifetime = dbconfig.connMaxLifetime
	}

	return nil
}
