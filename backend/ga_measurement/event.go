package ga_measurement

import (
	"context"
	"net/http"
	"strings"

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

func MustEventValues(ctx context.Context) *EventValues { // TODO: Move this func into interface/http package
	v, err := GetEventValues(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func getRemoteAddress(req *http.Request) string {
	xForwardedFor := req.Header.Get("X-Forwarded-For")
	if xForwardedFor == "" {
		return (strings.Split(req.RemoteAddr, ":"))[0]
	}
	return strings.TrimSpace((strings.Split(xForwardedFor, ","))[0])
}
