package http

import (
	"context"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/ga_measurement"
	"github.com/oinume/lekcije/backend/model"
	api_v1 "github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
)

type UserService struct {
	appLogger           *zap.Logger
	db                  *gorm.DB
	gaMeasurementClient ga_measurement.Client
}

func NewUserService(
	db *gorm.DB,
	appLogger *zap.Logger,
	gaMeasurementClient ga_measurement.Client,
) api_v1.User {
	return &UserService{
		appLogger:           appLogger,
		db:                  db,
		gaMeasurementClient: gaMeasurementClient,
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
	user, err := authenticateFromContext(ctx, s.db)
	if err != nil {
		return nil, err
	}
	return &api_v1.GetMeEmailResponse{Email: user.Email}, nil
}

func (s *UserService) UpdateMeEmail(
	ctx context.Context,
	request *api_v1.UpdateMeEmailRequest,
) (*api_v1.UpdateMeEmailResponse, error) {
	user, err := authenticateFromContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	// TODO: better validation
	email := request.Email
	if email == "" || !validateEmail(email) {
		return nil, twirp.InvalidArgumentError("email", "invalid email")
	}

	userService := model.NewUserService(s.db)
	if err := userService.UpdateEmail(user, request.Email); err != nil {
		return nil, err
	}

	go s.sendGAMeasurementEvent(
		ctx,
		ga_measurement.CategoryUser,
		"update",
		fmt.Sprint(user.ID),
		0,
		user.ID,
	)

	return &api_v1.UpdateMeEmailResponse{}, nil
}

func (s *UserService) UpdateMeNotificationTimeSpan(
	ctx context.Context,
	request *api_v1.UpdateMeNotificationTimeSpanRequest,
) (*api_v1.UpdateMeNotificationTimeSpanResponse, error) {
	user, err := authenticateFromContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	// TODO: validation
	timeSpanService := model.NewNotificationTimeSpanService(s.db)
	timeSpans := timeSpanService.NewNotificationTimeSpansFromPB(user.ID, request.NotificationTimeSpans)
	if err := timeSpanService.UpdateAll(user.ID, timeSpans); err != nil {
		return nil, err
	}

	go s.sendGAMeasurementEvent(
		ctx,
		ga_measurement.CategoryUser,
		"updateNotificationTimeSpan",
		fmt.Sprint(user.ID),
		0,
		user.ID,
	)

	return &api_v1.UpdateMeNotificationTimeSpanResponse{}, nil
}

func validateEmail(email string) bool {
	// TODO: better validation
	return strings.Contains(email, "@")
}

func (s *UserService) sendGAMeasurementEvent(
	ctx context.Context,
	category,
	action,
	label string,
	value int64,
	userID uint32,
) {
	err := s.gaMeasurementClient.SendEvent(
		ga_measurement.MustEventValues(ctx),
		category,
		action,
		fmt.Sprint(userID),
		0,
		userID,
	)
	if err != nil {
		s.appLogger.Warn("SendEvent() failed", zap.Error(err))
	}
}
