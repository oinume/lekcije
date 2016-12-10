package measurement

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Print

func TestEventParams_Values(t *testing.T) {
	a := assert.New(t)

	const (
		trackingID = "UA-2241989-19"
		clientID   = "123.456"
		category   = "category"
		action     = "action"
		label      = "label"
		value      = 100
	)
	params := NewEventParams("user-agent", trackingID, clientID, category, action)
	params.EventLabel = label
	params.EventValue = value
	v := params.Values()

	a.Equal(len(v), 9)
	a.Equal(trackingID, v.Get("tid"))
	a.Equal(clientID, v.Get("cid"))
	a.Equal(category, v.Get("ec"))
	a.Equal(action, v.Get("ea"))
	a.Equal(label, v.Get("el"))
	a.Equal(fmt.Sprint(value), v.Get("ev"))
}

type mockTransport struct {
	requestBody string
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("body = %v\n", string(body))
	t.requestBody = string(body)
	resp := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusOK,
		Status:     "OK",
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	return resp, nil
}

func TestClient_Do(t *testing.T) {
	a := assert.New(t)
	transport := &mockTransport{}
	httpClient := &http.Client{
		Transport: transport,
	}
	client := NewClient(httpClient)

	const (
		trackingID = "UA-2241989-19"
		clientID   = "123.456"
		category   = "category"
		action     = "action"
		label      = "label"
		value      = 100
	)
	params := NewEventParams("user-agent", trackingID, clientID, category, action)
	params.EventLabel = label
	params.EventValue = value

	err := client.Do(params)
	a.Nil(err)
	a.Contains(transport.requestBody, "cid=123.456") // TODO: Check other params
}
