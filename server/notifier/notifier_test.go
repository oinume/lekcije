package notifier

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/oinume/lekcije/server/emailer"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

var helper = model.NewTestHelper()
var _ = fmt.Print

type mockSenderTransport struct {
	sync.Mutex
	called int
}

func (t *mockSenderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.Lock()
	t.called++
	t.Unlock()
	time.Sleep(time.Millisecond * 500)
	resp := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusAccepted,
		Status:     "202 Accepted",
	}
	//resp.Header.Set("Content-Type", "text/html; charset=UTF-8")
	resp.Body = ioutil.NopCloser(strings.NewReader(""))
	return resp, nil
}

func TestMain(m *testing.M) {
	db := helper.DB()
	defer db.Close()
	helper.TruncateAllTables(db)
	os.Exit(m.Run())
}

func TestSendNotification(t *testing.T) {
	//a := assert.New(t)
	r := require.New(t)
	db := helper.DB()
	logger.InitializeAppLogger(os.Stdout, zapcore.DebugLevel)

	fetcherMockTransport, err := fetcher.NewMockTransport("../fetcher/testdata/5982.html")
	r.NoError(err)
	fetcherHTTPClient := &http.Client{
		Transport: fetcherMockTransport,
		Timeout:   5 * time.Second,
	}
	fetcher := fetcher.NewLessonFetcher(fetcherHTTPClient, 1, false, helper.LoadMCountries(), nil)

	var users []*model.User
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("oinume+%02d", i)
		user := helper.CreateUser(name, name+"@gmail.com")
		teacher := helper.CreateRandomTeacher()
		helper.CreateFollowingTeacher(user.ID, teacher)
		users = append(users, user)
	}

	senderHTTPClient := &http.Client{
		Transport: &mockSenderTransport{},
		Timeout:   5 * time.Second,
	}
	sender := emailer.NewSendGridSender(senderHTTPClient)
	n := NewNotifier(db, fetcher, true, sender)
	defer n.Close()

	for _, user := range users {
		err := n.SendNotification(user)
		r.NoError(err)
	}
}
