package usecase

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/randoms"
	"github.com/oinume/lekcije/backend/repository"
)

type UserAPIToken struct {
	userAPITokenRepo repository.UserAPIToken
}

func NewUserAPIToken(userAPITokenRepo repository.UserAPIToken) *UserAPIToken {
	return &UserAPIToken{
		userAPITokenRepo: userAPITokenRepo,
	}
}

func (u *UserAPIToken) Create(ctx context.Context, userID uint) (*model2.UserAPIToken, error) {
	// TODO: Idempotency
	uat := &model2.UserAPIToken{
		Token:  randoms.MustNewString(64),
		UserID: userID,
	}
	if err := u.userAPITokenRepo.Create(ctx, uat); err != nil {
		return nil, err
	}
	return uat, nil
}
