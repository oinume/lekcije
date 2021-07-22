package usecase

import (
	"context"

	model2 "github.com/oinume/lekcije/backend/model2c"
	"github.com/oinume/lekcije/backend/repository"
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
	userID uint32,
) error {
	return u.gaMeasurementRepo.SendEvent(ctx, values, category, action, label, value, userID)
}
