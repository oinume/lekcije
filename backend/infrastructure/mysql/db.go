package mysql

import (
	"context"
	"database/sql"

	"github.com/oinume/lekcije/backend/repository"
)

type dbRepository struct {
	db *sql.DB
}

func NewDB(db *sql.DB) repository.DB {
	return &dbRepository{db: db}
}

func (r *dbRepository) Transaction(ctx context.Context, f func(exec repository.Executor) error) error {
	return r.TransactionWithOptions(ctx, &sql.TxOptions{}, f)
}

func (r *dbRepository) TransactionWithOptions(ctx context.Context, opts *sql.TxOptions, f func(exec repository.Executor) error) error {
	tx, err := r.db.BeginTx(ctx, opts)
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
