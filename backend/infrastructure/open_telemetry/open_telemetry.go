package open_telemetry

import (
	"context"
	"io"
	"os"

	gcptrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	"github.com/oinume/lekcije/backend/domain/config"
)

type nopSpanExporter struct{}

func (e *nopSpanExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	return nil
}

func (e *nopSpanExporter) Shutdown(ctx context.Context) error {
	return nil
}

func NewTracerProvider(serviceName string, cfg *config.Vars) (*trace.TracerProvider, error) {
	var exporter trace.SpanExporter
	var err error
	switch cfg.Exporter {
	// TODO: jaeger
	// https://github.com/oinume/opencensus-client-trace-sample
	// https://github.com/open-telemetry/opentelemetry-go/blob/main/example/jaeger/main.go
	case "cloud_trace":
		exporter, err = gcptrace.New(gcptrace.WithProjectID(cfg.GCPProjectID))
		if err != nil {
			return nil, err
		}

		// Create trace provider with the exporter.
		//
		// By default it uses AlwaysSample() which samples all traces.
		// In a production environment or high QPS setup please use
		// probabilistic sampling.
		// Example:
		//   tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.0001)), ...)
		// defer provider.ForceFlush(ctx) // flushes any pending spans
		// otel.SetTracerProvider(provider)
	case "jaeger":
		exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Trace.ExporterURL)))
		if err != nil {
			return nil, err
		}
	case "stdout":
		exporter, err = NewStdoutExporter(os.Stdout)
		if err != nil {
			return nil, err
		}
	}

	if !cfg.Enable {
		exporter = &nopSpanExporter{}
	}

	r := NewResource(serviceName, config.DefaultVars.VersionHash, config.DefaultVars.ServiceEnv)
	return trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(r),
		trace.WithSampler(trace.TraceIDRatioBased(cfg.Trace.SamplingRate)),
	), nil
}

func NewStdoutExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

func NewResource(serviceName string, version string, environment string) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(version),
			attribute.String("environment", environment),
		),
	)
	return r
}
