package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/grpc_server"
	"github.com/oinume/lekcije/server/grpc_server/interceptor"
	"github.com/oinume/lekcije/server/route"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	bootstrap.CheckServerEnvVars()
}

func main() {
	port := config.ListenPort()
	grpcPort := config.GRPCListenPort()
	if port == grpcPort {
		log.Fatalf("Can't specify same port for a server.")
	}

	errors := make(chan error)
	go func() {
		errors <- startGRPCServer(grpcPort)
	}()

	go func() {
		errors <- startHTTPServer(grpcPort, port)
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
	grpc_server.RegisterAPIV1Server(server)
	// Register reflection service on gRPC server.
	reflection.Register(server)
	fmt.Printf("Starting gRPC server on %d\n", port)
	return server.Serve(lis)
}

func startHTTPServer(grpcPort, httpPort int) error {
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
	routes := route.Create(gatewayMux)

	fmt.Printf("Starting HTTP server on %v\n", httpPort)
	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), routes)
}
