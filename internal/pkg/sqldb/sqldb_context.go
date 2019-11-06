package sqldb

import (
	"context"
	"database/sql"
)

// GetContext function
func (db *DB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return db.follower.GetContext(ctx, dest, query, args...)
}

// SelectContext fuction
func (db *DB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return db.follower.SelectContext(ctx, dest, query, args...)
}

// QueryContext function
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.follower.QueryContext(ctx, query, args...)
}

// QueryRowContext function
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.follower.QueryRowContext(ctx, query, args...)
}

// ExecContext function
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.leader.ExecContext(ctx, query, args...)
}

// NamedExecContext function
func (db *DB) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return db.leader.NamedExecContext(ctx, query, arg)
}
