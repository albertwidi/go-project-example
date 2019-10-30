package sqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/albertwidi/go_project_example/internal/pkg/log/logger"
	"github.com/albertwidi/go_project_example/internal/pkg/log/logger/std"
	"github.com/jmoiron/sqlx"
)

// list of error
var (
	errConfigNil = errors.New("sqldb: config is nil")
)

// DB struct to hold all database connections
type DB struct {
	leader   *sqlx.DB
	follower *sqlx.DB

	logger logger.Logger
	config *Config
}

// Config of the db
type Config struct {
	Driver      string
	LeaderDSN   string
	FollowerDSN string
	Retry       int
	Logger      logger.Logger
}

// New sqldb object
func New(ctx context.Context, config *Config) (*DB, error) {
	var err error

	if config == nil {
		return nil, errConfigNil
	}

	if config.Logger == nil {
		config.Logger, err = std.New(nil)
		if err != nil {
			return nil, err
		}
	}

	db := DB{
		config: config,
	}
	leader, err := db.connectWithRetry(ctx, config.Driver, config.LeaderDSN, config.Retry)
	if err != nil {
		return nil, err
	}
	db.leader = leader

	follower, err := db.connectWithRetry(ctx, config.Driver, config.FollowerDSN, config.Retry)
	if err != nil {
		return nil, err
	}
	db.follower = follower

	return &db, nil
}

func (db *DB) connect(ctx context.Context, driver, dsn string) (*sqlx.DB, error) {
	sqlxdb, err := sqlx.ConnectContext(ctx, driver, dsn)
	if err != nil {
		return nil, err
	}
	return sqlxdb, err
}

func (db *DB) connectWithRetry(ctx context.Context, driver, dsn string, retry int) (*sqlx.DB, error) {
	var (
		sqlxdb *sqlx.DB
		err    error
	)

	if retry == 0 {
		sqlxdb, err = db.connect(ctx, driver, dsn)
		return nil, err
	}

	for x := 0; x < retry; x++ {
		sqlxdb, err = db.connect(ctx, driver, dsn)
		if err == nil {
			break
		}

		db.logger.Warnf("sqldb: failed to connect to %s with error %s", dsn, err.Error())
		db.logger.Warnf("sqldb: retrying to connect to %s. Retry: %d", dsn, x+1)

		if x+1 == retry && err != nil {
			db.logger.Errorf("sqldb: retry time exhausted, cannot connect to database: %s", err.Error())
			return nil, fmt.Errorf("sqldb: failed connect to database: %s", err.Error())
		}
		time.Sleep(time.Second * 3)
	}
	return sqlxdb, err
}

// Leader return leader database connection
func (db *DB) Leader() *sqlx.DB {
	return db.leader
}

// Follower return follower database connection
func (db *DB) Follower() *sqlx.DB {
	return db.follower
}

// SetMaxIdleConns to sql database
func (db *DB) SetMaxIdleConns(n int) {
	db.Leader().SetMaxIdleConns(n)
	db.Follower().SetMaxIdleConns(n)
}

// SetMaxOpenConns to sql database
func (db *DB) SetMaxOpenConns(n int) {
	db.Leader().SetMaxOpenConns(n)
	db.Follower().SetMaxOpenConns(n)
}

// SetConnMaxLifetime to sql database
func (db *DB) SetConnMaxLifetime(t time.Duration) {
	db.Leader().SetConnMaxLifetime(t)
	db.Follower().SetConnMaxLifetime(t)
}

// Get return one value in destination using relfection
func (db *DB) Get(dest interface{}, query string, args ...interface{}) error {
	return db.follower.Get(dest, query, args...)
}

// Select return more than one value in destintion using reflection
func (db *DB) Select(dest interface{}, query string, args ...interface{}) error {
	return db.follower.Select(dest, query, args...)
}

// Query function
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.follower.Query(query, args...)
}

// NamedQuery function
func (db *DB) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return db.follower.NamedQuery(query, arg)
}

// QueryRow function
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.follower.QueryRow(query, args...)
}

// Exec function
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.leader.Exec(query, args...)
}

// NamedExec execute query with named parameter
func (db *DB) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return db.leader.NamedExec(query, arg)
}

// Begin return sql transaction object, begin a transaction
func (db *DB) Begin() (*sql.Tx, error) {
	return db.leader.Begin()
}

// Beginx return sqlx transaction object, begin a transaction
func (db *DB) Beginx() (*sqlx.Tx, error) {
	return db.leader.Beginx()
}

// Rebind query
func (db *DB) Rebind(query string) string {
	return sqlx.Rebind(sqlx.BindType(db.config.Driver), query)
}

// Named return named query and parameters
func (db *DB) Named(query string, arg interface{}) (string, interface{}, error) {
	return sqlx.Named(query, arg)
}

// BindNamed return named query wrapped with bind
func (db *DB) BindNamed(query string, arg interface{}) (string, interface{}, error) {
	return sqlx.BindNamed(sqlx.BindType(db.config.Driver), query, arg)
}
