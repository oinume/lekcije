package repository

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
)

type UserAPIToken interface {
	Create(ctx context.Context, userAPIToken *model2.UserAPIToken) error
}
