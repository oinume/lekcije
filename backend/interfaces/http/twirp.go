package http

import (
	"context"

	api_v1 "github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
)

type UserService struct{}

func NewUserService() api_v1.User {
	return &UserService{}
}

func (s *UserService) Ping(
	ctx context.Context,
	request *api_v1.PingRequest,
) (*api_v1.PingResponse, error) {
	return &api_v1.PingResponse{}, nil
}

func (s *UserService) GetMe(
	ctx context.Context,
	request *api_v1.GetMeRequest,
) (*api_v1.GetMeResponse, error) {
	panic("implement me")
}

func (s *UserService) GetMeEmail(
	ctx context.Context,
	request *api_v1.GetMeEmailRequest,
) (*api_v1.GetMeEmailResponse, error) {
	panic("implement me")
}

func (s *UserService) UpdateMeEmail(
	ctx context.Context,
	request *api_v1.UpdateMeEmailRequest,
) (*api_v1.UpdateMeEmailResponse, error) {
	panic("implement me")
}

func (s *UserService) UpdateMeNotificationTimeSpan(
	ctx context.Context,
	request *api_v1.UpdateMeNotificationTimeSpanRequest,
) (*api_v1.UpdateMeNotificationTimeSpanResponse, error) {
	panic("implement me")
}
