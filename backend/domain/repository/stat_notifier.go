package repository

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
)

type StatNotifier interface {
	CreateOrUpdate(ctx context.Context, statNotifier *model2.StatNotifier) error
}
