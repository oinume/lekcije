package grpc_server

import (
	"fmt"
	"strings"

	api_v1 "github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
	"github.com/oinume/lekcije/server/context_data"
	"github.com/oinume/lekcije/server/event_logger"
	"github.com/oinume/lekcije/server/model"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type apiV1Server struct{}

func RegisterAPIV1Server(server *grpc.Server) {
	api_v1.RegisterAPIServer(server, &apiV1Server{})
}

func (s *apiV1Server) GetMe(
	ctx context.Context, in *api_v1.GetMeRequest,
) (*api_v1.GetMeResponse, error) {
	user, err := authorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	timeSpans := make([]*api_v1.NotificationTimeSpan, 0, 3)
	timeSpans = append(timeSpans, &api_v1.NotificationTimeSpan{
		FromHour:   12,
		FromMinute: 30,
		ToHour:     15,
		ToMinute:   00,
	})
	return &api_v1.GetMeResponse{
		Email: user.Email,
		NotificationTimeSpans: timeSpans,
	}, nil
}

func (s *apiV1Server) GetMeEmail(
	ctx context.Context, in *api_v1.GetMeEmailRequest,
) (*api_v1.GetMeEmailResponse, error) {
	user, err := authorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return &api_v1.GetMeEmailResponse{Email: user.Email}, nil
}

func (s *apiV1Server) UpdateMeEmail(
	ctx context.Context, in *api_v1.UpdateMeEmailRequest,
) (*api_v1.UpdateMeEmailResponse, error) {
	user, err := authorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: better validation
	email := in.Email
	if email == "" || !validateEmail(email) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Email")
	}

	userService := model.NewUserService(context_data.MustDB(ctx))
	if err := userService.UpdateEmail(user, in.Email); err != nil {
		return nil, err
	}
	go event_logger.SendGAMeasurementEvent2(event_logger.MustGAMeasurementEventValues(ctx), event_logger.CategoryUser, "update", fmt.Sprint(user.ID), 0, user.ID)

	return &api_v1.UpdateMeEmailResponse{}, nil
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

func validateEmail(email string) bool {
	// TODO: better validation
	return strings.Contains(email, "@")
}
