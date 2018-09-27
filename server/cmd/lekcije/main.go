package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/profiler"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/interfaces"
	interfaces_grpc "github.com/oinume/lekcije/server/interfaces/grpc"
	"github.com/oinume/lekcije/server/interfaces/grpc/interceptor"
	interfaces_http "github.com/oinume/lekcije/server/interfaces/http"
	"github.com/oinume/lekcije/server/interfaces/http/flash_message"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/util"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const maxDBConnections = 10

func main() {
	config.MustProcessDefault()
	port := config.DefaultVars.HTTPPort
	grpcPort := config.DefaultVars.GRPCPort
	if port == grpcPort {
		log.Fatalf("Can't specify same port for a server.")
	}

	if config.DefaultVars.EnableStackdriverProfiler {
		// TODO: Move to gcp package
		f, err := util.GenerateTempFileFromBase64String("", "gcloud-", config.DefaultVars.GcloudServiceKey)
		if err != nil {
			log.Fatalf("Failed to generate temp file: %v", err)
		}
		defer func() {
			os.Remove(f.Name())
		}()
		if err := profiler.Start(profiler.Config{
			ProjectID:      config.DefaultVars.GCPProjectID,
			Service:        "lekcije",
			ServiceVersion: "1.0.0", // TODO: release version?
			DebugLogging:   false,
		}, option.WithCredentialsFile(f.Name())); err != nil {
			log.Fatalf("Stackdriver profiler.Start failed: %v", err)
		}
	}

	db, err := model.OpenDB(
		config.DefaultVars.DBURL(),
		maxDBConnections,
		config.DefaultVars.DebugSQL,
	)
	if err != nil {
		log.Fatalf("model.OpenDB failed: %v", err)
	}
	defer db.Close()

	redis, err := model.OpenRedis(config.DefaultVars.RedisURL)
	if err != nil {
		log.Fatalf("model.OpenRedis failed: %v", err)
	}
	defer redis.Close()

	args := &interfaces.ServerArgs{
		DB:                db,
		FlashMessageStore: flash_message.NewStoreRedis(redis),
		Redis:             redis,
		SenderHTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	errors := make(chan error)
	go func() {
		errors <- startGRPCServer(grpcPort, args)
	}()

	go func() {
		errors <- startHTTPServer(grpcPort, port, args)
	}()

	for err := range errors {
		log.Fatal(err)
	}
}

func startGRPCServer(port int, args *interfaces.ServerArgs) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	opts = append(opts, interceptor.WithUnaryServerInterceptors())
	server := grpc.NewServer(opts...)
	interfaces_grpc.RegisterAPIV1Server(server, args) // TODO: RegisterAPIV1Service
	// Register reflection service on gRPC server.
	reflection.Register(server)
	fmt.Printf("Starting gRPC server on %d\n", port)
	return server.Serve(lis)
}

func startHTTPServer(grpcPort, httpPort int, args *interfaces.ServerArgs) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	muxOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		OrigName:     true,
		EmitDefaults: true,
	})
	gatewayMux := runtime.NewServeMux(muxOptions)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	endpoint := fmt.Sprintf("127.0.0.1:%d", grpcPort)
	if err := api_v1.RegisterAPIHandlerFromEndpoint(ctx, gatewayMux, endpoint, opts); err != nil {
		return err
	}
	server := interfaces_http.NewServer(args)
	mux := server.CreateRoutes(gatewayMux)
	fmt.Printf("Starting HTTP server on %v\n", httpPort)
	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), mux)
}
