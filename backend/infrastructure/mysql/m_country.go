package mysql

import (
	"context"
	"database/sql"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/model2"
)

type mCountryRepository struct {
	db *sql.DB
}

func NewMCountryRepository(db *sql.DB) repository.MCountry {
	return &mCountryRepository{
		db: db,
	}
}

func (r *mCountryRepository) FindAll(ctx context.Context) ([]*model2.MCountry, error) {
	return model2.MCountries().All(ctx, r.db)
}
