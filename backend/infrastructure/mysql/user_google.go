package mysql

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/model2"
)

type userGoogleRepository struct {
	db *sql.DB
}

func NewUserGoogleRepository(db *sql.DB) repository.UserGoogle {
	return &userGoogleRepository{db: db}
}

func (r *userGoogleRepository) CreateWithExec(ctx context.Context, exec repository.Executor, userGoogle *model2.UserGoogle) error {
	return userGoogle.Insert(ctx, exec, boil.Infer())
}

func (r *userGoogleRepository) DeleteByPKWithExec(ctx context.Context, exec repository.Executor, googleID string) error {
	ug := &model2.UserGoogle{GoogleID: googleID}
	if _, err := ug.Delete(ctx, exec); err != nil {
		return err
	}
	return nil
}

func (r *userGoogleRepository) FindByPKWithExec(ctx context.Context, exec repository.Executor, googleID string) (*model2.UserGoogle, error) {
	return model2.UserGoogles(model2.UserGoogleWhere.GoogleID.EQ(googleID)).One(ctx, exec)
}

func (r *userGoogleRepository) FindByUserIDWithExec(ctx context.Context, exec repository.Executor, userID uint) (*model2.UserGoogle, error) {
	return model2.UserGoogles(qm.Where("user_id = ?", userID)).One(ctx, exec)
}
