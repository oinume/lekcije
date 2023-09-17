package open_telemetry

import (
	"context"
	"fmt"
	"io"
	"os"

	gcptrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
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
	case "cloud_trace":
		exporter, err = gcptrace.New(gcptrace.WithProjectID(cfg.GCPProjectID))
		if err != nil {
			return nil, err
		}
	case "jaeger":
		client := otlptracehttp.NewClient()
		exporter, err = otlptrace.New(context.Background(), client)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
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
