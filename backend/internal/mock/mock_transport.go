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

func NewHTMLTransportFromString(content string) *HTMLTransport {
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

type ResponseTransport struct {
	sync.Mutex
	NumCalled    int
	responseFunc func(*ResponseTransport, *http.Request) *http.Response
}

func NewResponseTransport(f func(*ResponseTransport, *http.Request) *http.Response) *ResponseTransport {
	return &ResponseTransport{responseFunc: f}
}

func (t *ResponseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.Lock()
	t.NumCalled++
	t.Unlock()
	resp := t.responseFunc(t, req)
	if resp.StatusCode == 0 {
		resp.StatusCode = http.StatusOK
		resp.Status = http.StatusText(http.StatusOK)
	}
	if ct := resp.Header.Get("Content-Type"); ct == "" {
		resp.Header.Set("Content-Type", "text/plain; charset=UTF-8")
	}
	if resp.Body == nil {
		resp.Body = io.NopCloser(strings.NewReader(""))
	}
	return resp, nil
}
