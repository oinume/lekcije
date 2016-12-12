package measurement

// Google Analytics Measurement Protocol document
// https://developers.google.com/analytics/devguides/collection/protocol/v1/?hl=ja
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/oinume/lekcije/server/errors"
)

const (
	baseURL    = "https://www.google-analytics.com"
	collectURL = baseURL + "/collect"
	version    = 1
)

type Params interface {
	Validate() []error
	Values() url.Values
}

type CommonParams struct {
	version     int    // v
	trackingID  string // tid
	clientID    string // cid
	hitType     string // t
	AnonymousIP bool   // aip
	DataSource  string // ds
	UserID      string // uid
	UserAgent   string // ua
}

func (cp *CommonParams) Validate() []error {
	return nil
}

func (cp *CommonParams) Values() url.Values {
	v := url.Values{}
	v.Set("v", fmt.Sprint(cp.version))
	v.Set("tid", cp.trackingID)
	v.Set("cid", cp.clientID)
	v.Set("t", cp.hitType)
	if cp.AnonymousIP {
		v.Set("aip", "1")
	}
	if cp.DataSource != "" {
		v.Set("ds", cp.DataSource)
	}
	if cp.UserAgent != "" {
		v.Set("ua", cp.UserAgent)
	}
	if cp.UserID != "" {
		v.Set("uid", cp.UserID)
	}
	return v
}

func NewPageviewParams(userAgent, trackingID, clientID, documentHostname, page, title string) *PageviewParams {
	return &PageviewParams{
		CommonParams: &CommonParams{
			UserAgent:  userAgent,
			version:    version,
			trackingID: trackingID,
			clientID:   clientID,
			hitType:    "pageview",
		},
		documentHostname: documentHostname,
		page:             page,
		title:            title,
	}
}

type PageviewParams struct {
	*CommonParams
	documentHostname, page, title string
}

func (pp *PageviewParams) Validate() []error {
	return nil
}

func NewEventParams(userAgent, trackingID, clientID, eventCategory, eventAction string) *EventParams {
	return &EventParams{
		CommonParams: &CommonParams{
			UserAgent:  userAgent,
			version:    version,
			trackingID: trackingID,
			clientID:   clientID,
			hitType:    "event",
		},
		eventCategory: eventCategory,
		eventAction:   eventAction,
	}
}

type EventParams struct {
	*CommonParams
	eventCategory string
	eventAction   string
	EventLabel    string
	EventValue    int64
}

func (ep *EventParams) Validate() []error {
	// TODO:
	return nil
}

func (ep *EventParams) Values() url.Values {
	v := ep.CommonParams.Values()
	v.Set("ec", ep.eventCategory)
	v.Set("ea", ep.eventAction)
	if ep.EventLabel != "" {
		v.Set("el", ep.EventLabel)
	}
	if ep.EventValue != 0 {
		v.Set("ev", fmt.Sprint(ep.EventValue))
	}
	return v
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func (c *Client) Do(params Params) error {
	if errs := params.Validate(); len(errs) > 0 {
		message := ""
		for _, e := range errs {
			message += e.Error()
			message += ":"
		}
		return errors.Internalf("params.Validate() failed: %v", message)
	}

	v := params.Values()
	req, err := http.NewRequest("POST", collectURL, bytes.NewBufferString(v.Encode()))
	if err != nil {
		return errors.InternalWrapf(err, "http.NewRequest() failed")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.InternalWrapf(err, "httpClient.Do() failed: url=%v, values=%v", collectURL, v.Encode())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.InternalWrapf(err, "ioutil.ReadAll(resp.Body) failed")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Internalf("Call Measurement API failed: status=%v, body=%v", resp.StatusCode, string(body))
	}

	return nil
}

func GetClientID(cookie *http.Cookie) (string, error) {
	// Cookie value is like "GA1.2.1197953909.1480947778"
	values := strings.Split(cookie.Value, ".")
	if len(values) != 4 {
		return "", fmt.Errorf("The cookie is invalid format: value=", cookie.Value)
	}
	return strings.Join(values[2:], "."), nil
}

//func toURLValues(rv reflect.Value, values url.Values) url.Values {
//	for i := 0; i < rv.Type().NumField(); i++ {
//		fieldType := rv.Type().Field(i)
//		name := fieldType.Tag.Get("measurement")
//		values.Set(name, fmt.Sprint(rv.Field(i).Interface()))
//	}
//	return values
//}
