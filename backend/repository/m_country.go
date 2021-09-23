package repository

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
)

type MCountry interface {
	FindAll(ctx context.Context) ([]*model2.MCountry, error)
}
