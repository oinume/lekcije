package ga_measurement

// TODO: Define unit test in http
//func TestNewEventValuesFromRequest(t *testing.T) {
//	const (
//		userAgent  = "go-test"
//		ipOverride = "1.1.1.1"
//	)
//	tests := map[string]struct {
//		request *http.Request
//		want    *EventValues
//	}{
//		"normal": {
//			request: httptest.NewRequest("POST", "http://localhost/event", nil),
//			want: &EventValues{
//				UserAgentOverride: userAgent,
//				ClientID:          "",
//				DocumentHostName:  "localhost",
//				DocumentPath:      "/event",
//				DocumentTitle:     "/event",
//				DocumentReferrer:  "",
//				IPOverride:        ipOverride,
//			},
//		},
//	}
//	for name, test := range tests {
//		test.request.Header.Set("User-Agent", userAgent)
//		test.request.Header.Set("X-Forwarded-For", ipOverride)
//		t.Run(name, func(t *testing.T) {
//			if got := NewEventValuesFromRequest(test.request); !reflect.DeepEqual(got, test.want) {
//				t.Errorf("NewEventValuesFromRequest() = %v, want %v", got, test.want)
//			}
//		})
//	}
//}
