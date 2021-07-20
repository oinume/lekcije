package http

import (
	"context"
	"net/http"

	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/errors"
	model2 "github.com/oinume/lekcije/backend/model2c"
)

type eventValuesKey struct{}

func NewEventValuesFromRequest(req *http.Request) *model2.GAMeasurementEvent {
	// Ignore if client id is not set
	clientID, _ := context_data.GetTrackingID(req.Context())
	return &model2.GAMeasurementEvent{
		UserAgentOverride: req.UserAgent(),
		ClientID:          clientID,
		DocumentHostName:  req.Host,
		DocumentPath:      req.URL.Path,
		DocumentTitle:     req.URL.Path,
		DocumentReferrer:  req.Referer(),
		IPOverride:        getRemoteAddress(req),
	}
}

func WithEventValues(ctx context.Context, v *model2.GAMeasurementEvent) context.Context {
	return context.WithValue(ctx, eventValuesKey{}, v)
}

func GetEventValues(ctx context.Context) (*model2.GAMeasurementEvent, error) {
	v := ctx.Value(eventValuesKey{})
	if value, ok := v.(*model2.GAMeasurementEvent); ok {
		return value, nil
	} else {
		return nil, errors.NewInternalError(
			errors.WithMessage("failed get value from context"),
		)
	}
}

func MustEventValues(ctx context.Context) *model2.GAMeasurementEvent {
	v, err := GetEventValues(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
