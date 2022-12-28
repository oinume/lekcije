package mysql

import (
	"context"
	"database/sql"

	"github.com/oinume/lekcije/backend/domain/repository"
)

type dbRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) repository.DB {
	return &dbRepository{db: db}
}

func (r *dbRepository) Transaction(ctx context.Context, f func(exec repository.Executor) error) error {
	return r.TransactionWithOptions(ctx, &sql.TxOptions{}, f)
}

func (r *dbRepository) TransactionWithOptions(ctx context.Context, opts *sql.TxOptions, f func(exec repository.Executor) error) error {
	return transactionWithOptions(ctx, r.db, opts, f)
}

func transaction(ctx context.Context, db *sql.DB, f func(exec repository.Executor) error) error {
	return transactionWithOptions(ctx, db, &sql.TxOptions{}, f)
}

func transactionWithOptions(
	ctx context.Context,
	db *sql.DB,
	opts *sql.TxOptions,
	f func(exec repository.Executor) error,
) error {
	tx, err := db.BeginTx(ctx, opts)
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
