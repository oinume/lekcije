package http

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	model2 "github.com/oinume/lekcije/backend/model2c"
)

func TestNewGAMeasurementEventFromRequest(t *testing.T) {
	const (
		userAgent  = "go-test"
		ipOverride = "1.1.1.1"
	)
	tests := map[string]struct {
		request *http.Request
		want    *model2.GAMeasurementEvent
	}{
		"normal": {
			request: httptest.NewRequest("POST", "http://localhost/event", nil),
			want: &model2.GAMeasurementEvent{
				UserAgentOverride: userAgent,
				ClientID:          "",
				DocumentHostName:  "localhost",
				DocumentPath:      "/event",
				DocumentTitle:     "/event",
				DocumentReferrer:  "",
				IPOverride:        ipOverride,
			},
		},
	}
	for name, test := range tests {
		test.request.Header.Set("User-Agent", userAgent)
		test.request.Header.Set("X-Forwarded-For", ipOverride)
		t.Run(name, func(t *testing.T) {
			if got := newGAMeasurementEventFromRequest(test.request); !reflect.DeepEqual(got, test.want) {
				t.Errorf("newGAMeasurementEventFromRequest() = %v, want %v", got, test.want)
			}
		})
	}
}
