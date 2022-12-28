package usecase

import (
	"context"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/model2"
)

type StatNotifier struct {
	statNotifierRepo repository.StatNotifier
}

func NewStatNotifier(statNotifierRepo repository.StatNotifier) *StatNotifier {
	return &StatNotifier{
		statNotifierRepo: statNotifierRepo,
	}
}

func (u *StatNotifier) CreateOrUpdate(ctx context.Context, statNotifier *model2.StatNotifier) error {
	return u.statNotifierRepo.CreateOrUpdate(ctx, statNotifier)
}
