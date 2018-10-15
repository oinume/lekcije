package gcp

import (
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/profiler"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/util"
	"go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"google.golang.org/api/option"
)

func StartStackdriverProfiler(c *config.Vars, service, serviceVersion string) error {
	f, err := util.GenerateTempFileFromBase64String("", "gcloud-", c.GcloudServiceKey)
	if err != nil {
		log.Fatalf("Failed to generate temp file: %v", err)
	}
	defer func() {
		os.Remove(f.Name())
	}()
	if err := profiler.Start(profiler.Config{
		ProjectID:      config.DefaultVars.GCPProjectID,
		Service:        service,
		ServiceVersion: serviceVersion, // TODO: release version?
		DebugLogging:   true,
	}, option.WithCredentialsFile(f.Name())); err != nil {
		return err
	}
	return nil
}

func EnableStackdriverTrace(c *config.Vars) error {
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: c.GCPProjectID,
	})
	if err != nil {
		return err
		//log.Fatalf("stackdriver.NewExporter failed: %v", err)
	}
	//exporter := &exporter.PrintExporter{}
	// Export to Stackdriver Monitoring.
	view.RegisterExporter(exporter)

	// Subscribe views to see stats in Stackdriver Monitoring.
	if err := view.Register(
		ochttp.ClientLatencyView,
		ochttp.ClientResponseBytesView,
	); err != nil {
		//log.Fatalf("view.Register failed: %v", err)
		return err
	}

	// Export to Stackdriver Trace.
	trace.RegisterExporter(exporter)

	// Automatically add a Stackdriver trace header to outgoing requests:
	http.DefaultClient.Transport = &ochttp.Transport{
		Propagation: &tracecontext.HTTPFormat{},
	}

	return nil
}
