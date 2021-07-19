package usecase

import (
	"context"

	"github.com/oinume/lekcije/backend/ga_measurement"
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
	values *ga_measurement.EventValues,
	category, action, label string,
	value int64,
	userID uint32,
) error {
	return u.gaMeasurementRepo.SendEvent(ctx, values, category, action, label, value, userID)
}
