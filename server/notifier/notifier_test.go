package notifier

import (
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/stretchr/testify/assert"
)

var helper = model.NewTestHelper()

type mockTransport struct {
	sync.Mutex
	called int
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
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

func TestMain(m *testing.M) {
	db := helper.DB()
	defer db.Close()
	helper.TruncateAllTables(db)
	os.Exit(m.Run())
}

func TestSendNotification(t *testing.T) {
	a := assert.New(t)
	db := helper.DB()
	logger.InitializeAppLogger(os.Stdout)

	client := &http.Client{
		Transport: &mockTransport{},
		Timeout:   5 * time.Second,
	}
	fetcher := fetcher.NewTeacherLessonFetcher(client, 1, false, helper.LoadMCountries(), nil)
	n := NewNotifier(db, fetcher, true, false)

	user := helper.CreateUser("oinume", "oinume@gmail.com")
	// TODO: follow teacher
	err := n.SendNotification(user)
	a.Nil(err)
}
