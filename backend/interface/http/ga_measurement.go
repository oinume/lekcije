package http

import (
	"net/http"

	"github.com/oinume/lekcije/backend/context_data"
	model2 "github.com/oinume/lekcije/backend/model2c"
)

func newGAMeasurementEventFromRequest(req *http.Request) *model2.GAMeasurementEvent {
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
