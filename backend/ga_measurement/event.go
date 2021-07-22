package ga_measurement

// TODO: Remove this file

import (
	"context"

	"github.com/oinume/lekcije/backend/errors"
	model2 "github.com/oinume/lekcije/backend/model2c"
)

const (
	CategoryEmail            = "email"
	CategoryUser             = "user"
	CategoryFollowingTeacher = "followingTeacher"
)

type eventValuesKey struct{}

type EventValues = model2.GAMeasurementEvent

func GetEventValues(ctx context.Context) (*EventValues, error) {
	v := ctx.Value(eventValuesKey{})
	if value, ok := v.(*EventValues); ok {
		return value, nil
	} else {
		return nil, errors.NewInternalError(
			errors.WithMessage("failed get value from context"),
		)
	}
}

func MustEventValues(ctx context.Context) *EventValues {
	v, err := GetEventValues(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
