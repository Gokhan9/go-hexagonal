package repository

import (
	"context"
	"database/sql"
)

type TxKey struct {
}

func WithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, TxKey{}, tx)
}

func GetTx(ctx context.Context) *sql.Tx {

	if tx, ok := ctx.Value(TxKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}
