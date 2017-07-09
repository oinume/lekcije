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
)

var helper = model.NewTestHelper()
var _ = fmt.Print

type mockFetcherTransport struct {
	sync.Mutex
	called int
}

func (t *mockFetcherTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.Lock()
	t.called++
	t.Unlock()
	resp := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusOK,
		Status:     "200 OK",
	}
	resp.Header.Set("Content-Type", "text/html; charset=UTF-8")

	// TODO: file location
	file, err := os.Open("../fetcher/testdata/5982.html")
	if err != nil {
		return nil, err
	}
	resp.Body = file // Close() will be called by client
	return resp, nil
}

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
	logger.InitializeAppLogger(os.Stdout)

	fetcherHTTPClient := &http.Client{
		Transport: &mockFetcherTransport{},
		Timeout:   5 * time.Second,
	}
	fetcher := fetcher.NewTeacherLessonFetcher(fetcherHTTPClient, 1, false, helper.LoadMCountries(), nil)

	usersData := []struct {
		name  string
		email string
	}{
		{"oinume", "oinume@gmail.com"},
		{"oinume2", "oinume+2@gmail.com"},
		{"oinume3", "oinume+3@gmail.com"},
	}
	var users []*model.User
	for _, u := range usersData {
		user := helper.CreateUser(u.name, u.email)
		teacher := helper.CreateRandomTeacher()
		helper.CreateFollowingTeacher(user.ID, teacher)
		users = append(users, user)
	}

	senderHTTPClient := &http.Client{
		Transport: &mockSenderTransport{},
		Timeout:   5 * time.Second,
	}
	sender := emailer.NewSendGridSender(senderHTTPClient)
	//sender := &emailer.NoSender{}
	n := NewNotifier(db, fetcher, true, sender)

	for _, user := range users {
		err := n.SendNotification(user)
		r.Nil(err)
	}
}
