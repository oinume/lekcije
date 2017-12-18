package grpc_server

import (
	"fmt"

	api_v1 "github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"github.com/oinume/lekcije/server/logger"
	"go.uber.org/zap"
)

type apiV1Server struct{}

func RegisterAPIV1Server(server *grpc.Server) {
	api_v1.RegisterAPIServer(server, &apiV1Server{})
}

func (s *apiV1Server) GetMeEmail(
	ctx context.Context, in *api_v1.GetMeEmailRequest,
) (*api_v1.GetMeEmailResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no metadata")
	}
	if v, ok := md["api-token"]; ok {
		logger.App.Info("api-token", zap.String("api-token", v[0]))
	}
	// TODO: implement
	return &api_v1.GetMeEmailResponse{Email:"a@foo.com"}, nil
}
