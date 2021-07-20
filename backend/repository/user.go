package repository

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
)

type User interface {
	CreateWithExec(ctx context.Context, exec Executor, user *model2.User) error
	FindByGoogleID(ctx context.Context, googleID string) (*model2.User, error)
	FindByGoogleIDWithExec(ctx context.Context, exec Executor, googleID string) (*model2.User, error)
	UpdateEmail(ctx context.Context, id uint, email string) error
}
