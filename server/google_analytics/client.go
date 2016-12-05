package google_analytics

// Google Analytics Measurement Protocol document
// https://developers.google.com/analytics/devguides/collection/protocol/v1/?hl=ja
import (
	"fmt"
	"net/http"
	"strings"
	"net/url"
)

const (
	baseURL = "http://www.google-analytics.com"
	version = 1
)

type Params interface {
	//Validate() []error
	Values() url.Values
}

type CommonParams struct {
	version int
	trackingID string
	clientID string
	hitType string
}

type PageviewParams struct {
	*CommonParams
	documentHostname, page, title string
}

type EventParams struct {
	*CommonParams
	EventCategory, EventAction, EventLabel, EventValue string
}

func GetClientID(cookie *http.Cookie) (string, error) {
	// Cookie value is like "GA1.2.1197953909.1480947778"
	values := strings.Split(cookie.Value, ".")
	if len(values) != 4 {
		return "", fmt.Errorf("The cookie is invalid format: value=", cookie.Value)
	}
	return strings.Join(values[2:], "."), nil
}

func NewPageviewParams(trackingID, clientID, documentHostname, page, title string) (*PageviewParams) {
	return &PageviewParams{
		CommonParams: &CommonParams{
			version: version,
			trackingID: trackingID,
			clientID: clientID,
			hitType: "pageview",
		},
		documentHostname: documentHostname,
		page: page,
		title: title,
	}
}

func NewEventParams(trackingID, clientID, eventCategory, eventAction, eventLabel, eventValue string) (*EventParams) {
	if trackingID == "" {
		panic("trackingID is required.") // TODO: return error?
	}
	return &EventParams{
		CommonParams: &CommonParams{
			version: version,
			trackingID: trackingID,
			clientID: clientID,
			hitType: "event",
		},
		EventCategory: eventCategory,
		EventAction: eventAction,
		EventLabel: eventLabel,
		EventValue: eventValue,
	}
}

func (cp *CommonParams) Values() url.Values {
	v := url.Values{}
	return v
}

func (ep *EventParams) Values() url.Values {
	v := ep.CommonParams.Values()
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
	//values := params.Values()
	return nil
}
