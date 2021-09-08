package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/profiler"
	"github.com/rollbar/rollbar-go"
	"go.opencensus.io/trace"
	"goji.io/v3"

	"github.com/oinume/lekcije/backend/config"
	"github.com/oinume/lekcije/backend/di"
	"github.com/oinume/lekcije/backend/event_logger"
	"github.com/oinume/lekcije/backend/gcp"
	"github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
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
	configVars := config.DefaultVars
	port := configVars.HTTPPort

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
		credential, err := gcp.WithCredentialsJSONFromBase64String(configVars.GCPServiceAccountKey)
		if err != nil {
			log.Fatalf("WithCredentialsJSONFromBase64String failed: %v", err)
		}

		// TODO: Move to gcp package
		if err := profiler.Start(profiler.Config{
			ProjectID:      configVars.GCPProjectID,
			Service:        serviceName,
			ServiceVersion: "1.0.0", // TODO: release version?
			DebugLogging:   false,
		}, credential); err != nil {
			log.Fatalf("Stackdriver profiler.Start failed: %v", err)
		}
	}

	gormDB, err := model.OpenDB(
		configVars.DBURL(),
		maxDBConnections,
		configVars.DebugSQL,
	)
	if err != nil {
		log.Fatalf("model.OpenDB failed: %v", err)
	}
	defer gormDB.Close()

	accessLogger := logger.NewAccessLogger(os.Stdout)
	appLogger := logger.NewAppLogger(os.Stderr, logger.NewLevel("info")) // TODO: flag
	rollbarClient := rollbar.New(
		configVars.RollbarAccessToken,
		configVars.ServiceEnv,
		configVars.VersionHash,
		"", "/",
	)
	defer rollbarClient.Close()

	args := &interfaces.ServerArgs{
		AccessLogger:        accessLogger,
		AppLogger:           appLogger,
		DB:                  gormDB.DB(),
		FlashMessageStore:   flash_message.NewStoreMySQL(gormDB),
		GAMeasurementClient: ga_measurement.NewClient(nil, event_logger.New(accessLogger)),
		GormDB:              gormDB,
		RollbarClient:       rollbarClient,
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
	errorRecorder := di.NewErrorRecorderUsecase(args.AppLogger, args.RollbarClient)
	server := interfaces_http.NewServer(args, errorRecorder)
	oauthServer := di.NewOAuthServer(args.AppLogger, args.DB, args.GAMeasurementClient, args.RollbarClient, args.SenderHTTPClient)
	errorRecorderHooks := interfaces_http.NewErrorRecorderHooks(errorRecorder)
	meServer := di.NewMeServer(
		args.AppLogger, args.GormDB, errorRecorderHooks, args.GAMeasurementClient,
	)

	mux := goji.NewMux()
	server.Setup(mux)
	oauthServer.Setup(mux)
	interfaces_http.SetupTwirpServer(mux, meServer)

	fmt.Printf("Starting HTTP server on %v\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
