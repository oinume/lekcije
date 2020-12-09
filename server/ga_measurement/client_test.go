package ga_measurement

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/oinume/lekcije/server/event_logger"
)

type fakeTransport struct{}

func (t *fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	panic("fakeTransport.RoundTrip")
	//b, err := httputil.DumpRequest(req, true)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(b))
	//resp := &http.Response{
	//	Header:     make(http.Header),
	//	Request:    req,
	//	StatusCode: http.StatusOK,
	//	Status:     "200 OK",
	//}
	//resp.Header.Set("Content-Type", "text/html; charset=UTF-8")
	//resp.Body = ioutil.NopCloser(strings.NewReader("ok"))
	//return resp, nil
}

type fakeErrorTransport struct{}

func (t *fakeErrorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	//b, err := httputil.DumpRequest(req, true)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(b))
	resp := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusInternalServerError,
		Status:     "500 Internal Server Error",
	}
	resp.Header.Set("Content-Type", "text/html; charset=UTF-8")
	resp.Body = ioutil.NopCloser(strings.NewReader("error"))
	return resp, nil
}

func Test_client_SendEvent(t *testing.T) {
	type args struct {
		values   *EventValues
		category string
		action   string
		label    string
		value    int64
		userID   uint32
	}
	tests := map[string]struct {
		transport http.RoundTripper
		args      args
		wantErr   bool
	}{
		"normal": {
			transport: &fakeTransport{},
			args: args{
				values: &EventValues{
					UserAgentOverride: "go test",
					//ClientID:          "",
					DocumentHostName: "www.lekcije.com",
					DocumentPath:     "/url",
					DocumentTitle:    "hoge",
					DocumentReferrer: "",
					IPOverride:       "192.168.99.100",
				},
			},
			wantErr: false,
		},
		"error_from_ga": {
			transport: &fakeErrorTransport{},
			args: args{
				values: &EventValues{
					UserAgentOverride: "go test",
					//ClientID:          "",
					DocumentHostName: "www.lekcije.com",
					DocumentPath:     "/url",
					DocumentTitle:    "hoge",
					DocumentReferrer: "",
					IPOverride:       "192.168.99.100",
				},
			},
			wantErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := NewClient(
				&http.Client{Transport: test.transport},
				zap.NewNop(),
				event_logger.New(zap.NewNop()),
			)

			err := c.SendEvent(
				test.args.values,
				test.args.category,
				test.args.action,
				test.args.label,
				test.args.value,
				test.args.userID,
			)
			if (err != nil) != test.wantErr {
				t.Fatalf("unexpected SendEvent result: wantErr=%v, err=%v", test.wantErr, err)
			}
		})
	}
}
