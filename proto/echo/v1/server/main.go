package main

// curl -X POST -d '{"value":"hoge"}' http://localhost:50052/echo

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/oinume/lekcije/proto-gen/go/proto/echo/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s *server) Echo(ctx context.Context, in *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{Message: "Echo:" + in.Message}, nil
}

func (s *server) EchoV2(ctx context.Context, in *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{Message: "Echo:" + in.Message}, nil
}

const (
	defaultGRPCPort    = "50051"
	defaultGatewayPort = "50052"
)

func startGRPCServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	echo.RegisterEchoServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	fmt.Println("Starting gRPC server on " + port)
	return s.Serve(lis)
}

func startGatewayServer(grpcPort, port string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gatewayMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := echo.RegisterEchoHandlerFromEndpoint(ctx, gatewayMux, "127.0.0.1:"+grpcPort, opts); err != nil {
		return err
	}

	//mux := http.NewServeMux()
	//mux.Handle("/v1/", gatewayMux)
	//mux.Handle("/", http.FileServer(http.Dir("swagger-ui")))

	fmt.Println("Starting gateway server on " + port)
	return http.ListenAndServe(":"+port, gatewayMux)
}

//func cors(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		Set(w, AccessControl{
//			Origin:         "*",
//			AllowedMethods: []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "DELETE", "PATCH"},
//		})
//		next.ServeHTTP(w, r)
//	})
//}

func main() {
	grpcPort := os.Getenv("GPRC_PORT")
	gatewayPort := os.Getenv("PORT")
	if grpcPort == "" {
		grpcPort = defaultGRPCPort
	}
	if gatewayPort == "" {
		gatewayPort = defaultGatewayPort
	}
	if grpcPort == gatewayPort {
		log.Fatalf("Can't specify same port.")
	}

	errors := make(chan error)
	go func() {
		errors <- startGRPCServer(grpcPort)
	}()

	go func() {
		errors <- startGatewayServer(grpcPort, gatewayPort)
	}()

	for err := range errors {
		log.Fatal(err)
	}
}
