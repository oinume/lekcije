package open_census

import (
	"fmt"
	"log"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/zipkin"
	open_zipkin "github.com/openzipkin/zipkin-go"
	zipkin_http "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/trace"
	"google.golang.org/api/option"

	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/gcp"
)

type FlushFunc func()

type nopExporter struct{}

func (e *nopExporter) ExportSpan(s *trace.SpanData) {}

func NewExporter(c *config.Vars, service string, alwaysSample bool) (trace.Exporter, FlushFunc, error) {
	if !c.EnableTrace {
		return &nopExporter{}, func() {}, nil
	}

	var exporter trace.Exporter
	var flush FlushFunc
	if c.ZipkinReporterURL == "" {
		if c.GCPProjectID == "" {
			return nil, nil, fmt.Errorf("no exporter configuration")
		}
		credential, err := gcp.WithCredentialsJSONFromBase64String(c.GCPServiceAccountKey)
		if err != nil {
			return nil, nil, err
		}
		e, err := stackdriver.NewExporter(stackdriver.Options{
			ProjectID: c.GCPProjectID,
			// MetricPrefix helps uniquely identify your metrics.
			MetricPrefix:            service,
			TraceClientOptions:      []option.ClientOption{credential},
			MonitoringClientOptions: []option.ClientOption{credential},
		})
		if err != nil {
			log.Fatalf("Failed to create the Stackdriver exporter: %v", err)
		}

		exporter = e
		// It is imperative to invoke flush before your main function exits
		flush = e.Flush
		if alwaysSample {
			trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
		} else {
			trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(0.1)})
		}
	} else {
		// 1. Configure exporter to export traces to Zipkin.
		localEndpoint, err := open_zipkin.NewEndpoint(service, "192.168.1.5:5454")
		if err != nil {
			return nil, func() {}, err
		}
		reporter := zipkin_http.NewReporter(c.ZipkinReporterURL)
		e := zipkin.NewExporter(reporter, localEndpoint)
		exporter = e
		flush = func() {
			_ = reporter.Close()
		}
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	return exporter, flush, nil
}
