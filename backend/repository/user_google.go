package repository

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
)

type UserGoogle interface {
	CreateWithExec(ctx context.Context, exec Executor, userGoogle *model2.UserGoogle) error
	DeleteByPKWithExec(ctx context.Context, exec Executor, googleID string) error
	FindByPKWithExec(ctx context.Context, exec Executor, googleID string) (*model2.UserGoogle, error)
	FindByUserIDWithExec(ctx context.Context, exec Executor, userID uint) (*model2.UserGoogle, error)
}
