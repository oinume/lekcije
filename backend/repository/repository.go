package repository

import (
	"context"
	"database/sql"
)

// Executor can perform SQL queries. It's for abstraction of database/sql.DB
type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type TxBeginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Begin() (*sql.Tx, error)
}

func Transaction(ctx context.Context, txb TxBeginner, f func(tx Executor) error) error {
	return TransactionWithOptions(ctx, txb, &sql.TxOptions{}, f)
}

func TransactionWithOptions(ctx context.Context, txb TxBeginner, opts *sql.TxOptions, f func(tx Executor) error) error {
	tx, err := txb.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	if err := f(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

type DB interface {
	Transaction(ctx context.Context, f func(exec Executor) error) error
	TransactionWithOptions(ctx context.Context, opts *sql.TxOptions, f func(exec Executor) error) error
}
