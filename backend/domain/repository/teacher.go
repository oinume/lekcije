package repository

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
)

type Teacher interface {
	Create(ctx context.Context, teacher *model2.Teacher) error
	CreateOrUpdate(ctx context.Context, teacher *model2.Teacher) error
	FindByID(ctx context.Context, id uint) (*model2.Teacher, error)
	FindByIDs(ctx context.Context, ids []uint) ([]*model2.Teacher, error)
	IncrementFetchErrorCount(ctx context.Context, id uint, value int) error
}
