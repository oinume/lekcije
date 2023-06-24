package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/profiler"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
	"github.com/rollbar/rollbar-go"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"goji.io/v3"
	"goji.io/v3/pat"

	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/event_logger"
	"github.com/oinume/lekcije/backend/infrastructure/firebase"
	"github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
	"github.com/oinume/lekcije/backend/infrastructure/gcp"
	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	"github.com/oinume/lekcije/backend/infrastructure/open_telemetry"
	interfaces "github.com/oinume/lekcije/backend/interface"
	"github.com/oinume/lekcije/backend/interface/graphql"
	"github.com/oinume/lekcije/backend/interface/graphql/generated"
	"github.com/oinume/lekcije/backend/interface/graphql/resolver"
	interfaces_http "github.com/oinume/lekcije/backend/interface/http"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/registry"
)

const (
	maxDBConnections = 10
	serviceName      = "lekcije"
)

func main() {
	ctx := context.Background()
	gae := os.Getenv("GAE_APPLICATION")
	if gae != "" {
		for _, env := range []string{"MYSQL_USER", "MYSQL_PASSWORD", "MYSQL_HOST"} {
			if os.Getenv(env) == "" {
				log.Fatalf("Environment variable %q is required", env)
			}
			if err := berglas.Replace(ctx, env); err != nil {
				log.Fatalf("berglas.Replace failed: %v", err)
			}
		}
	}
	config.MustProcessDefault()
	configVars := config.DefaultVars
	port := configVars.HTTPPort

	tracerProvider, err := open_telemetry.NewTracerProvider("server", config.DefaultVars)
	if err != nil {
		log.Fatalf("NewTraceProvider failed: %v", err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tracerProvider)

	firebaseApp, err := firebase.NewApp(ctx, configVars.FirebaseProjectID, configVars.FirebaseServiceAccountKeyBase64)
	if err != nil {
		log.Fatalf("firebase.NewApp failed: %v\n", err)
	}
	firebaseAuthClient, err := firebaseApp.Auth(ctx)
	if err != nil {
		log.Fatalf("firebaseApp.Auth failed: %v\n", err)
	}

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
		FirebaseAuthClient:  firebaseAuthClient,
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
	mux := goji.NewMux()

	// TODO: graceful shutdown
	errorRecorder := registry.NewErrorRecorderUsecase(args.AppLogger, args.RollbarClient)
	userAPIToken := registry.NewUserAPITokenUsecase(args.DB)
	server := interfaces_http.NewServer(args, errorRecorder, userAPIToken)
	oauthServer := registry.NewOAuthServer(args.AppLogger, args.DB, args.FirebaseAuthClient, args.GAMeasurementClient, args.RollbarClient, args.SenderHTTPClient)
	mCountryList := registry.MustNewMCountryList(context.Background(), args.DB)
	server.Setup(mux)
	oauthServer.Setup(mux)

	gqlResolver := resolver.NewResolver(
		mysql.NewFollowingTeacherRepository(args.DB),
		registry.NewFollowingTeacherUsecase(args.AppLogger, args.DB, mCountryList),
		registry.NewGAMeasurementUsecase(args.GAMeasurementClient),
		mysql.NewNotificationTimeSpanRepository(args.DB),
		registry.NewNotificationTimeSpanUsecase(args.DB),
		mysql.NewTeacherRepository(args.DB),
		mysql.NewUserRepository(args.DB),
		registry.NewUserUsecase(args.DB, args.FirebaseAuthClient),
	)
	gqlSchema := generated.NewExecutableSchema(generated.Config{
		Resolvers: gqlResolver,
	})
	gqlServer := handler.NewDefaultServer(gqlSchema)
	gqlMiddleware := graphql.NewMiddleware(args.AppLogger)
	gqlServer.AroundOperations(gqlMiddleware.AroundOperations)
	gqlServer.SetErrorPresenter(gqlMiddleware.ErrorPresenter)

	const gqlPath = "/graphql"
	mux.Handle(pat.Get(gqlPath), gqlServer)
	mux.Handle(pat.Post(gqlPath), gqlServer)
	if !config.DefaultVars.IsProductionEnv() {
		mux.Handle(pat.Get("/playground"), playground.Handler("GraphQL playground", gqlPath))
	}

	otelHandler := otelhttp.NewHandler(
		mux,
		"server",
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)
	fmt.Printf("Starting HTTP server on %v\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), otelHandler)
}
