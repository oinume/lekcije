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
	db := helper.DB()
	logger.InitializeAppLogger(os.Stdout, zapcore.DebugLevel)

	fetcherMockTransport, err := fetcher.NewMockTransport("../fetcher/testdata/5982.html")
	if err != nil {
		t.Fatalf("fetcher.NewMockTransport failed: err=%v", err)
	}
	fetcherHTTPClient := &http.Client{
		Transport: fetcherMockTransport,
		Timeout:   5 * time.Second,
	}
	fetcher := fetcher.NewLessonFetcher(fetcherHTTPClient, 1, false, helper.LoadMCountries(), nil)

	var users []*model.User
	const numOfUsers = 10
	for i := 0; i < numOfUsers; i++ {
		name := fmt.Sprintf("oinume+%02d", i)
		user := helper.CreateUser(name, name+"@gmail.com")
		teacher := helper.CreateRandomTeacher()
		helper.CreateFollowingTeacher(user.ID, teacher)
		users = append(users, user)
	}

	senderTransport := &mockSenderTransport{}
	senderHTTPClient := &http.Client{
		Transport: senderTransport,
		Timeout:   5 * time.Second,
	}
	sender := emailer.NewSendGridSender(senderHTTPClient)
	n := NewNotifier(db, fetcher, true, sender)
	defer n.Close()

	for _, user := range users {
		err := n.SendNotification(user)
		if err != nil {
			t.Fatalf("SendNotification failed: err=%v", err)
		}
	}
	if got, want := senderTransport.called, numOfUsers; got != want {
		t.Errorf("unexpected senderTransport.called: got=%v, want=%v", got, want)
	}
}
