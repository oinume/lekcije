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

func NewExporter(c *config.Vars, service string) (trace.Exporter, error) {
	var exporter trace.Exporter
	if c.ZipkinReporterURL == "" {
		if c.GCPProjectID == "" {
			return nil, fmt.Errorf("no exporter configuration")
		}

		sd, err := stackdriver.NewExporter(stackdriver.Options{
			ProjectID: c.GCPProjectID,
			// MetricPrefix helps uniquely identify your metrics.
			MetricPrefix: service,
		})
		if err != nil {
			log.Fatalf("Failed to create the Stackdriver exporter: %v", err)
		}
		// It is imperative to invoke flush before your main function exits
		defer sd.Flush() // TODO: outside function

		// Register it as a trace exporter
		trace.RegisterExporter(sd) // TODO: outside function
		exporter = sd
	} else {
		// 1. Configure exporter to export traces to Zipkin.
		localEndpoint, err := open_zipkin.NewEndpoint(service, "192.168.1.5:5454")
		if err != nil {
			return nil, err
		}
		reporter := zipkin_http.NewReporter(c.ZipkinReporterURL)
		ze := zipkin.NewExporter(reporter, localEndpoint)
		trace.RegisterExporter(ze) // TODO: outside function

		// 2. Configure 100% sample rate, otherwise, few traces will be sampled.
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

		exporter = ze
	}

	return exporter, nil
}
