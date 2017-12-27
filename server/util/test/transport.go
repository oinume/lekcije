package test

import (
	"net/http"
	"os"
	"sync"
)

type MockFetcherTransport struct {
	sync.Mutex
	NumCalled int
	mockHTML  string
}

func NewMockFetcherTransport(mockHTML string) *MockFetcherTransport {
	return &MockFetcherTransport{
		mockHTML: mockHTML,
	}
}

func (t *MockFetcherTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.Lock()
	t.NumCalled++
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
	resp.Body = file // Close() will be NumCalled by client
	return resp, nil
}
