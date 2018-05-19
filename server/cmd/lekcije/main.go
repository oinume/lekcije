package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/interfaces"
	interfaces_grpc "github.com/oinume/lekcije/server/interfaces/grpc"
	"github.com/oinume/lekcije/server/interfaces/grpc/interceptor"
	interfaces_http "github.com/oinume/lekcije/server/interfaces/http"
	"github.com/oinume/lekcije/server/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const maxDBConnections = 5

func main() {
	config.MustProcessDefault()
	port := config.DefaultVars.HTTPPort
	grpcPort := config.DefaultVars.GRPCPort
	if port == grpcPort {
		log.Fatalf("Can't specify same port for a server.")
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
	//ctx = context_data.SetDB(ctx, db)

	redisClient, err := model.OpenRedis(config.DefaultVars.RedisURL)
	if err != nil {
		log.Fatalf("model.OpenRedis failed: %v", err)
	}
	defer redisClient.Close()
	//_, c = flash_message.NewStoreRedisAndSetToContext(c, redisClient)

	args := &interfaces.ServerArgs{
		DB:          db,
		RedisClient: redisClient,
	}
	errors := make(chan error)
	go func() {
		errors <- startGRPCServer(grpcPort)
	}()

	go func() {
		errors <- startHTTPServer(grpcPort, port, args)
	}()

	for err := range errors {
		log.Fatal(err)
	}
}

func startGRPCServer(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	opts = append(opts, interceptor.WithUnaryServerInterceptors())
	server := grpc.NewServer(opts...)
	interfaces_grpc.RegisterAPIV1Server(server)
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
