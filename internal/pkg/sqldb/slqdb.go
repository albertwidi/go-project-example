package sqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// list of error
var (
	errConfigNil = errors.New("sqldb: config is nil")
)

// DB struct to hold all database connections
type DB struct {
	driver   string
	leader   *sqlx.DB
	follower *sqlx.DB
}

// Wrap leader and follower sqlx object to one DB object
// this is for easier usage, so user doesn't have to specify leader or follower
// all exec is going to leader, all query is going to follower
func Wrap(ctx context.Context, leader, follower *sqlx.DB) (*DB, error) {
	if leader.DriverName() != follower.DriverName() {
		return nil, fmt.Errorf("sqldb: leader and follower driver is not matched. leader = %s follower = %s", leader.DriverName(), follower.DriverName())
	}

	db := DB{
		driver:   leader.DriverName(),
		leader:   leader,
		follower: follower,
	}
	return &db, nil
}

// ConnectOptions to list options when connect to the db
type ConnectOptions struct {
	Retry                 int
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
}

// Connect to a new database
func Connect(ctx context.Context, driver, dsn string, connOpts *ConnectOptions) (*sqlx.DB, error) {
	opts := connOpts
	if opts == nil {
		opts = &ConnectOptions{}
	}

	db, err := connectWithRetry(ctx, driver, dsn, opts.Retry)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(opts.MaxOpenConnections)
	db.SetMaxIdleConns(opts.MaxIdleConnections)
	db.SetConnMaxLifetime(opts.ConnectionMaxLifetime)
	return db, nil
}

func connectWithRetry(ctx context.Context, driver, dsn string, retry int) (*sqlx.DB, error) {
	var (
		sqlxdb *sqlx.DB
		err    error
	)

	if retry == 0 {
		sqlxdb, err = sqlx.ConnectContext(ctx, driver, dsn)
		if err != nil {
			return nil, err
		}
		return sqlxdb, err
	}

	for x := 0; x < retry; x++ {
		sqlxdb, err = sqlx.ConnectContext(ctx, driver, dsn)
		if err == nil {
			break
		}
		if x+1 == retry && err != nil {
			return nil, fmt.Errorf("sqldb: failed connect to database: %s", err.Error())
		}
		time.Sleep(time.Second * 3)
	}
	return sqlxdb, err
}

// Close all database connection to leader and replica
func (db *DB) Close() error {
	if err := db.leader.Close(); err != nil {
		return err
	}
	if err := db.follower.Close(); err != nil {
		return err
	}
	return nil
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
	return sqlx.Rebind(sqlx.BindType(db.driver), query)
}

// Named return named query and parameters
func (db *DB) Named(query string, arg interface{}) (string, interface{}, error) {
	return sqlx.Named(query, arg)
}

// BindNamed return named query wrapped with bind
func (db *DB) BindNamed(query string, arg interface{}) (string, interface{}, error) {
	return sqlx.BindNamed(sqlx.BindType(db.driver), query, arg)
}
