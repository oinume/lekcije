package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"reflect"
	"testing"
	"time"

	"go.opencensus.io/trace"
)

type fakeExporter struct {
	spanNames []string
}

func (e *fakeExporter) ExportSpan(s *trace.SpanData) {
	e.spanNames = append(e.spanNames, s.Name)
}

func TestNewHTTPClientTracer(t *testing.T) {
	exporter := &fakeExporter{}
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	trace.RegisterExporter(exporter)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer ts.Close()

	ctx := context.Background()
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}
	tracer := NewHTTPClientTracer(
		ctx,
		"test.",
		[]trace.Attribute{trace.StringAttribute("url", ts.URL)},
		fmt.Sprintf("url:%s", ts.URL),
	)
	req = req.WithContext(httptrace.WithClientTrace(ctx, tracer.Trace()))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Do failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	want := []string{"test.connect", "test.getConn", "test.waitForResponse"}
	if got := exporter.spanNames; !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected span names: got = %+v, want = %+v", got, want)
	}
}
