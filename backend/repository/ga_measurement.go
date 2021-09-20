package repository

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
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
