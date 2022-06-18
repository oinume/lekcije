package mysql

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/model2"
)

type userAPITokenRepository struct {
	db *sql.DB
}

func NewUserAPITokenRepository(db *sql.DB) repository.UserAPIToken {
	return &userAPITokenRepository{
		db: db,
	}
}

func (r *userAPITokenRepository) Create(ctx context.Context, userAPIToken *model2.UserAPIToken) error {
	return userAPIToken.Insert(ctx, r.db, boil.Infer())
}

func (r *userAPITokenRepository) DeleteByUserIDAndToken(ctx context.Context, userID uint, token string) error {
	userAPIToken := &model2.UserAPIToken{
		UserID: userID,
		Token:  token,
	}
	_, err := userAPIToken.Delete(ctx, r.db)
	return err
}
