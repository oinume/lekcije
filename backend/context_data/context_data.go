package context_data

import (
	"context"

	"github.com/morikuni/failure"

	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
)

type (
	apiTokenKey     struct{}
	eventValuesKey  struct{}
	loggedInUserKey struct{}
	trackingIDKey   struct{}
)

func GetLoggedInUser(ctx context.Context) (*model.User, error) {
	value := ctx.Value(loggedInUserKey{})
	if user, ok := value.(*model.User); ok {
		return user, nil
	}
	return nil, errors.NewNotFoundError(errors.WithMessage("Failed to get loggedInUser from context"))
}

func MustLoggedInUser(ctx context.Context) *model.User {
	user, err := GetLoggedInUser(ctx)
	if err != nil {
		panic(err)
	}
	return user
}

func SetLoggedInUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, loggedInUserKey{}, user)
}

func SetTrackingID(ctx context.Context, trackingID string) context.Context {
	return context.WithValue(ctx, trackingIDKey{}, trackingID)
}

func GetTrackingID(ctx context.Context) (string, error) {
	value := ctx.Value(trackingIDKey{})
	if trackingID, ok := value.(string); ok {
		return trackingID, nil
	}
	return "", errors.NewNotFoundError(errors.WithMessage("Failed to get trackingID from context"))
}

func MustTrackingID(ctx context.Context) string {
	trackingID, err := GetTrackingID(ctx)
	if err != nil {
		panic(err)
	}
	return trackingID
}

func SetAPIToken(ctx context.Context, apiToken string) context.Context {
	return context.WithValue(ctx, apiTokenKey{}, apiToken)
}

func GetAPIToken(ctx context.Context) (string, error) {
	value := ctx.Value(apiTokenKey{})
	if apiToken, ok := value.(string); ok {
		return apiToken, nil
	}
	return "", errors.NewNotFoundError(errors.WithMessage("failed to get api token from context"))
}

func SetGAMeasurementEvent(ctx context.Context, v *model2.GAMeasurementEvent) context.Context {
	return context.WithValue(ctx, eventValuesKey{}, v)
}

func GetGAMeasurementEvent(ctx context.Context) (*model2.GAMeasurementEvent, error) {
	v := ctx.Value(eventValuesKey{})
	if value, ok := v.(*model2.GAMeasurementEvent); ok {
		return value, nil
	} else {
		return nil, failure.New(errors.Internal, failure.Message("failed get value from context"))
	}
}

func MustGAMeasurementEvent(ctx context.Context) *model2.GAMeasurementEvent {
	v, err := GetGAMeasurementEvent(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
