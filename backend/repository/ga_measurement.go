package repository

import (
	"context"

	"github.com/oinume/lekcije/backend/ga_measurement"
)

type GAMeasurement interface {
	SendEvent(
		ctx context.Context,
		values *ga_measurement.EventValues,
		category,
		action,
		label string,
		value int64,
		userID uint32,
	) error
}
