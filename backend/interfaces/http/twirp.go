package http

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/oinume/lekcije/backend/model"
	api_v1 "github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) api_v1.User {
	return &UserService{
		db: db,
	}
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
	user, err := authenticateFromContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	timeSpansService := model.NewNotificationTimeSpanService(s.db)
	timeSpans, err := timeSpansService.FindByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	timeSpansPB, err := timeSpansService.NewNotificationTimeSpansPB(timeSpans)
	if err != nil {
		return nil, err
	}

	mPlan, err := model.NewMPlanService(s.db).FindByPK(user.PlanID)
	if err != nil {
		return nil, err
	}

	return &api_v1.GetMeResponse{
		UserId:                int32(user.ID),
		Email:                 user.Email,
		NotificationTimeSpans: timeSpansPB,
		MPlan: &api_v1.MPlan{
			Id:   int32(mPlan.ID),
			Name: mPlan.Name,
		},
	}, nil
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
