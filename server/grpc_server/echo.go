package grpc_server

import (
	"context"

	"github.com/oinume/lekcije/proto-gen/go/proto/echo/v1"
	"google.golang.org/grpc"
)

type echoServer struct{}

func RegisterEchoServer(server *grpc.Server) {
	echo.RegisterEchoServer(server, &echoServer{})
}

func (s *echoServer) Echo(ctx context.Context, in *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{Message: "Echo:" + in.Message}, nil
}
func (s *echoServer) EchoV2(ctx context.Context, in *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{Message: "EchoV2:" + in.Message}, nil
}
