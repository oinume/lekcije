package registration_email

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/randoms"
)

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

	user := modeltest.NewUser(func(u *model2.User) {
		u.ID = uint(randoms.MustNewInt64(10000000))
	})
	transport := &mockSenderTransport{}
	httpClient := &http.Client{
		Transport: transport,
	}

	appLogger := logger.NewAppLogger(os.Stdout, zapcore.InfoLevel)
	sender := NewEmailSender(httpClient, appLogger)
	err := sender.Send(context.Background(), user)
	r.NoError(err)
	a.Equal(1, transport.called)
	// TODO: assert content
}
