package repository

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
)

type GAMeasurement interface { // TODO: Rename to ActivityRepository
	SendEvent( // TODO: Rename to CreateActivity
		ctx context.Context,
		values *model2.GAMeasurementEvent, // TODO: Rename to ActivityEvent
		category,
		action,
		label string,
		value int64,
		userID uint32,
	) error
}
