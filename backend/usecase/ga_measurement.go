package usecase

import (
	"context"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/model2"
)

type GAMeasurement struct {
	gaMeasurementRepo repository.GAMeasurement
}

func NewGAMeasurement(gaMeasurementRepo repository.GAMeasurement) *GAMeasurement {
	return &GAMeasurement{
		gaMeasurementRepo: gaMeasurementRepo,
	}
}

func (u *GAMeasurement) SendEvent(
	ctx context.Context,
	values *model2.GAMeasurementEvent,
	category, action, label string,
	value int64,
	userID uint32, // TODO: uint
) error {
	return u.gaMeasurementRepo.SendEvent(ctx, values, category, action, label, value, userID)
}
