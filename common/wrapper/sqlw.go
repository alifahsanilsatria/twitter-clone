package wrapper

import (
	"context"
	"database/sql"
)

type SQLWrapper interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type SQLTxWrapper interface {
	Rollback() error
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Commit() error
}
