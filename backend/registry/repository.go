package registry

import (
	"context"
	"database/sql"

	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	"github.com/oinume/lekcije/backend/model2"
)

func MustNewMCountryList(ctx context.Context, db *sql.DB) *model2.MCountryList {
	mCountries, err := mysql.NewMCountryRepository(db).FindAll(ctx)
	if err != nil {
		panic(err)
	}
	return model2.NewMCountryList(mCountries)
}
