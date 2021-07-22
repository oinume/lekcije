package mysql

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/repository"
)

type userGoogleRepository struct {
	db *sql.DB
}

func NewUserGoogleRepository(db *sql.DB) repository.UserGoogle {
	return &userGoogleRepository{db: db}
}

func (r *userGoogleRepository) FindByUserIDWithExec(ctx context.Context, exec repository.Executor, userID uint) (*model2.UserGoogle, error) {
	ug, err := model2.UserGoogles(qm.Where("user_id = ?", userID)).One(ctx, exec)
	if err != nil {
		return nil, err
	}
	return ug, nil
}

func (r *userGoogleRepository) CreateWithExec(ctx context.Context, exec repository.Executor, userGoogle *model2.UserGoogle) error {
	return userGoogle.Insert(ctx, exec, boil.Infer())
}
