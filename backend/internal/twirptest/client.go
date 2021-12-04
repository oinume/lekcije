package twirptest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type JSONError struct {
	Code string            `json:"code"`
	Msg  string            `json:"msg"`
	Meta map[string]string `json:"meta,omitempty"`
}

func (e *JSONError) Error() string {
	return fmt.Sprintf("Code=%s, Msg=%s, Meta=%v", e.Code, e.Msg, e.Meta)
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
	) (int, *JSONError)
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
) (int, error) {
	t.Helper()

	marshaler := protojson.MarshalOptions{
		UseProtoNames: true,
	}
	body, err := marshaler.Marshal(request)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	req, err := http.NewRequest("POST", path, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		var je JSONError
		if err := json.NewDecoder(resp.Body).Decode(&je); err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
		return resp.StatusCode, &je
	}

	rawRespBody := json.RawMessage{}
	if err := json.NewDecoder(resp.Body).Decode(&rawRespBody); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	if err = unmarshaler.Unmarshal(rawRespBody, response); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	return resp.StatusCode, nil
}
