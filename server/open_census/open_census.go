package open_census

import (
	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/oinume/lekcije/server/config"
	open_zipkin "github.com/openzipkin/zipkin-go"
	zipkin_http "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/trace"
)

func RegisterExporter(c *config.Vars) error {
	if c.ZipkinReporterURL == "" {
		if c.GCPProjectID == "" {
			return nil
		}

	} else {
		// 1. Configure exporter to export traces to Zipkin.
		localEndpoint, err := open_zipkin.NewEndpoint("go-quickstart", "192.168.1.5:5454")
		if err != nil {
			return err
		}
		reporter := zipkin_http.NewReporter(c.ZipkinReporterURL)
		ze := zipkin.NewExporter(reporter, localEndpoint)
		trace.RegisterExporter(ze)

		// 2. Configure 100% sample rate, otherwise, few traces will be sampled.
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	return nil
}
