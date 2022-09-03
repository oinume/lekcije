package mock

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type HTMLTransport struct {
	sync.Mutex
	NumCalled int
	content   string
}

func NewHTMLTransport(path string) (*HTMLTransport, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("os.Open failed: path=%v, err=%v", path, err)
	}
	b, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read file failed: err=%v", err)
	}
	return &HTMLTransport{
		content: string(b),
	}, nil
}

func NewMockTransportFromHTML(content string) *HTMLTransport {
	return &HTMLTransport{
		content: content,
	}
}

func (t *HTMLTransport) RoundTrip(req *http.Request) (*http.Response, error) {
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
	resp.Body = io.NopCloser(strings.NewReader(t.content))
	return resp, nil
}
