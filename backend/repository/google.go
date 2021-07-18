package repository

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
)

type UserGoogle interface {
	FindByUserIDWithExec(ctx context.Context, exec Executor, userID uint) (*model2.UserGoogle, error)
	CreateWithExec(ctx context.Context, exec Executor, userGoogle *model2.UserGoogle) error
}
