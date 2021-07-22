package repository

import (
	"context"

	model2 "github.com/oinume/lekcije/backend/model2c"
)

type GAMeasurement interface {
	SendEvent(
		ctx context.Context,
		values *model2.GAMeasurementEvent,
		category,
		action,
		label string,
		value int64,
		userID uint32,
	) error
}
