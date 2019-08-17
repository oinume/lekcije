package ga_measurement

import (
	"context"
	"net/http"
	"strings"

	"github.com/oinume/lekcije/server/context_data"
	"github.com/oinume/lekcije/server/errors"
)

type eventValuesKey struct{}

type EventValues struct {
	UserAgentOverride string
	ClientID          string
	DocumentHostName  string
	DocumentPath      string
	DocumentTitle     string
	DocumentReferrer  string
	IPOverride        string
}

func NewEventValuesFromRequest(req *http.Request) *EventValues {
	// Ignore if client id is not set
	clientID, _ := context_data.GetTrackingID(req.Context())
	return &EventValues{
		UserAgentOverride: req.UserAgent(),
		ClientID:          clientID,
		DocumentHostName:  req.Host,
		DocumentPath:      req.URL.Path,
		DocumentTitle:     req.URL.Path,
		DocumentReferrer:  req.Referer(),
		IPOverride:        getRemoteAddress(req),
	}
}

func WithEventValues(ctx context.Context, v *EventValues) context.Context {
	return context.WithValue(ctx, eventValuesKey{}, v)
}

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

func getRemoteAddress(req *http.Request) string {
	xForwardedFor := req.Header.Get("X-Forwarded-For")
	if xForwardedFor == "" {
		return (strings.Split(req.RemoteAddr, ":"))[0]
	}
	return strings.TrimSpace((strings.Split(xForwardedFor, ","))[0])
}
