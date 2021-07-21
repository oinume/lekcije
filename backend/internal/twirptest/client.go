package twirptest

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/protobuf/jsonpb" //nolint:staticcheck
	"github.com/golang/protobuf/proto"  //nolint:staticcheck
)

type HTTPError struct {
	StatusCode int
	BodyString string
	BodyBytes  []byte
}

func (e *HTTPError) Error() string {
	if e.BodyString != "" {
		return fmt.Sprintf("StatusCode=%d, BodyString=%q", e.StatusCode, e.BodyString)
	} else {
		return fmt.Sprintf("StatusCode=%d, BodyBytes=%q", e.StatusCode, e.BodyBytes)
	}
}

type Client interface {
	SendRequest(
		ctx context.Context,
		t *testing.T,
		handler http.Handler,
		path string,
		request proto.Message,
		response proto.Message,
		wantStatusCode int,
	)
}

type JSONClient struct{}

func NewJSONClient() *JSONClient {
	return &JSONClient{}
}

func (jc *JSONClient) SendRequest(
	ctx context.Context,
	t *testing.T,
	handler http.Handler,
	path string,
	request proto.Message,
	response proto.Message,
	wantStatusCode int,
) error {
	t.Helper()

	var body bytes.Buffer
	marshaler := &jsonpb.Marshaler{}
	if err := marshaler.Marshal(&body, request); err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	req, err := http.NewRequest("POST", path, &body)
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != wantStatusCode {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("ReadAll failed: %v", err)
		}
		return &HTTPError{
			StatusCode: resp.StatusCode,
			BodyString: string(b),
		}
	}

	unmarshaler := &jsonpb.Unmarshaler{}
	if err := unmarshaler.Unmarshal(resp.Body, response); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	return nil
}
