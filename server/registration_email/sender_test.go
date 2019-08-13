package registration_email

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/oinume/lekcije/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var helper = model.NewTestHelper()
var _ = fmt.Print

type mockSenderTransport struct {
	sync.Mutex
	called      int
	requestBody string
}

func (t *mockSenderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.Lock()
	t.called++
	defer t.Unlock()
	time.Sleep(time.Millisecond * 500)
	resp := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusAccepted,
		Status:     "202 Accepted",
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return resp, err
	}
	t.requestBody = string(body)
	defer req.Body.Close()
	resp.Body = ioutil.NopCloser(strings.NewReader(""))
	return resp, nil
}

func TestEmailSender_Send(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	user := helper.CreateRandomUser(t)
	transport := &mockSenderTransport{}
	httpClient := &http.Client{
		Transport: transport,
	}

	sender := NewEmailSender(httpClient)
	err := sender.Send(context.Background(), user)
	r.NoError(err)
	a.Equal(1, transport.called)
	// TODO: assert content
}
