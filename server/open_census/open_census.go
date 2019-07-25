package open_census

import (
	"fmt"
	"log"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/oinume/lekcije/server/config"
	open_zipkin "github.com/openzipkin/zipkin-go"
	zipkin_http "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/trace"
)

type FlushFunc func()

func NewExporter(c *config.Vars, service string) (trace.Exporter, FlushFunc, error) {
	var exporter trace.Exporter
	var flush FlushFunc

	if c.ZipkinReporterURL == "" {
		if c.GCPProjectID == "" {
			return nil, func() {}, fmt.Errorf("no exporter configuration")
		}

		sd, err := stackdriver.NewExporter(stackdriver.Options{
			ProjectID: c.GCPProjectID,
			// MetricPrefix helps uniquely identify your metrics.
			MetricPrefix: service,
		})
		if err != nil {
			log.Fatalf("Failed to create the Stackdriver exporter: %v", err)
		}

		exporter = sd
		// It is imperative to invoke flush before your main function exits
		flush = sd.Flush
	} else {
		// 1. Configure exporter to export traces to Zipkin.
		localEndpoint, err := open_zipkin.NewEndpoint(service, "192.168.1.5:5454")
		if err != nil {
			return nil, func() {}, err
		}
		reporter := zipkin_http.NewReporter(c.ZipkinReporterURL)
		ze := zipkin.NewExporter(reporter, localEndpoint)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

		exporter = ze
		flush = func() {
			_ = reporter.Close()
		}
	}

	return exporter, flush, nil
}
