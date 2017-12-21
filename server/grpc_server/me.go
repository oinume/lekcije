package grpc_server

import (
	api_v1 "github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/oinume/lekcije/server/context_data"
	"github.com/oinume/lekcije/server/model"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type apiV1Server struct{}

func RegisterAPIV1Server(server *grpc.Server) {
	api_v1.RegisterAPIServer(server, &apiV1Server{})
}

func (s *apiV1Server) GetMeEmail(
	ctx context.Context, in *api_v1.GetMeEmailRequest,
) (*api_v1.GetMeEmailResponse, error) {
	user, err := authorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return &api_v1.GetMeEmailResponse{Email:user.Email}, nil
}

func authorizeFromContext(ctx context.Context) (*model.User, error) {
	apiToken, err := context_data.GetAPIToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "No api token found")
	}
	userService := model.NewUserService(context_data.MustDB(ctx))
	user, err := userService.FindLoggedInUser(apiToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "No user found")
	}
	return user, nil
}
