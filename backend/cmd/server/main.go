package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/profiler"
	"go.opencensus.io/trace"
	"goji.io/v3"

	"github.com/oinume/lekcije/backend/config"
	"github.com/oinume/lekcije/backend/di"
	"github.com/oinume/lekcije/backend/event_logger"
	"github.com/oinume/lekcije/backend/ga_measurement"
	"github.com/oinume/lekcije/backend/gcp"
	interfaces "github.com/oinume/lekcije/backend/interface"
	interfaces_http "github.com/oinume/lekcije/backend/interface/http"
	"github.com/oinume/lekcije/backend/interface/http/flash_message"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/open_census"
)

const (
	maxDBConnections = 10
	serviceName      = "lekcije"
)

func main() {
	config.MustProcessDefault()
	port := config.DefaultVars.HTTPPort

	exporter, flush, err := open_census.NewExporter(
		config.DefaultVars,
		serviceName,
		!config.DefaultVars.IsProductionEnv(),
	)
	if err != nil {
		log.Fatalf("NewExporter failed: %v", err)
	}
	defer flush()
	trace.RegisterExporter(exporter)

	if config.DefaultVars.EnableStackdriverProfiler {
		credential, err := gcp.WithCredentialsJSONFromBase64String(config.DefaultVars.GCPServiceAccountKey)
		if err != nil {
			log.Fatalf("WithCredentialsJSONFromBase64String failed: %v", err)
		}

		// TODO: Move to gcp package
		if err := profiler.Start(profiler.Config{
			ProjectID:      config.DefaultVars.GCPProjectID,
			Service:        serviceName,
			ServiceVersion: "1.0.0", // TODO: release version?
			DebugLogging:   false,
		}, credential); err != nil {
			log.Fatalf("Stackdriver profiler.Start failed: %v", err)
		}
	}

	gormDB, err := model.OpenDB(
		config.DefaultVars.DBURL(),
		maxDBConnections,
		config.DefaultVars.DebugSQL,
	)
	if err != nil {
		log.Fatalf("model.OpenDB failed: %v", err)
	}
	defer gormDB.Close()

	accessLogger := logger.NewAccessLogger(os.Stdout)
	appLogger := logger.NewAppLogger(os.Stderr, logger.NewLevel("info")) // TODO: flag
	args := &interfaces.ServerArgs{
		AccessLogger:      accessLogger,
		AppLogger:         appLogger,
		DB:                gormDB.DB(),
		FlashMessageStore: flash_message.NewStoreMySQL(gormDB),
		GAMeasurementClient: ga_measurement.NewClient(
			nil,
			appLogger,
			event_logger.New(accessLogger),
		),
		GormDB: gormDB,
		SenderHTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	errors := make(chan error)
	go func() {
		errors <- startHTTPServer(port, args)
	}()

	for err := range errors {
		log.Fatal(err)
	}
}

func startHTTPServer(port int, args *interfaces.ServerArgs) error {
	// TODO: graceful shutdown
	server := interfaces_http.NewServer(args)
	oauthServer := di.NewOAuthServer(args.AppLogger, args.DB, args.GAMeasurementClient, args.SenderHTTPClient)
	userServer := di.NewUserServer(
		args.AppLogger, args.GormDB, args.GAMeasurementClient,
	)

	mux := goji.NewMux()
	server.Setup(mux)
	oauthServer.Setup(mux)
	interfaces_http.SetupTwirpServer(mux, userServer)

	fmt.Printf("Starting HTTP server on %v\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
