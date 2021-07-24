package http

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/jinzhu/gorm"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/context_data"
	model2 "github.com/oinume/lekcije/backend/model2c"
	"github.com/oinume/lekcije/backend/usecase"
	api_v1 "github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
)

type UserService struct {
	appLogger                   *zap.Logger
	db                          *gorm.DB
	gaMeasurementUsecase        *usecase.GAMeasurement
	notificationTimeSpanUsecase *usecase.NotificationTimeSpan
	userUsecase                 *usecase.User
}

func NewUserService(
	db *gorm.DB, // TODO: Remove
	appLogger *zap.Logger,
	gaMeasurementUsecase *usecase.GAMeasurement,
	notificationTimeSpanUsecase *usecase.NotificationTimeSpan,
	userUsecase *usecase.User,
) api_v1.User {
	return &UserService{
		appLogger:                   appLogger,
		db:                          db,
		gaMeasurementUsecase:        gaMeasurementUsecase,
		notificationTimeSpanUsecase: notificationTimeSpanUsecase,
		userUsecase:                 userUsecase,
	}
}

func (s *UserService) Ping(
	_ context.Context,
	_ *api_v1.PingRequest,
) (*api_v1.PingResponse, error) {
	return &api_v1.PingResponse{}, nil
}

func (s *UserService) GetMe(
	ctx context.Context,
	_ *api_v1.GetMeRequest,
) (*api_v1.GetMeResponse, error) {
	user, err := authenticateFromContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	timeSpans, err := s.notificationTimeSpanUsecase.FindByUserID(ctx, uint(user.ID))
	if err != nil {
		return nil, err
	}
	timeSpansProto, err := NotificationTimeSpansProto(timeSpans)
	if err != nil {
		return nil, err
	}

	return &api_v1.GetMeResponse{
		UserId:                int32(user.ID),
		Email:                 user.Email,
		NotificationTimeSpans: timeSpansProto,
		MPlan:                 nil, // TODO: Remove
	}, nil
}

func (s *UserService) GetMeEmail(
	ctx context.Context,
	_ *api_v1.GetMeEmailRequest,
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

	email := request.Email
	if email == "" || !validateEmail(email) {
		return nil, twirp.InvalidArgumentError("email", "invalid email")
	}
	duplicate, err := s.userUsecase.IsDuplicateEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if duplicate {
		return nil, twirp.InvalidArgumentError("email", "email exists")
	}

	if err := s.userUsecase.UpdateEmail(ctx, uint(user.ID), email); err != nil {
		return nil, err
	}

	go func() {
		if err := s.gaMeasurementUsecase.SendEvent(
			ctx,
			context_data.MustGAMeasurementEvent(ctx),
			model2.GAMeasurementEventCategoryUser,
			"update",
			fmt.Sprint(user.ID),
			0,
			user.ID,
		); err != nil {
			panic(err) // TODO: Better error handling
		}
	}()

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
	userID := uint(user.ID)
	timeSpans := NotificationTimeSpansModel(request.NotificationTimeSpans, userID)
	if err := s.notificationTimeSpanUsecase.UpdateAll(ctx, userID, timeSpans); err != nil {
		return nil, err
	}

	go func() {
		if err := s.gaMeasurementUsecase.SendEvent(
			ctx,
			context_data.MustGAMeasurementEvent(ctx),
			model2.GAMeasurementEventCategoryUser,
			"updateNotificationTimeSpan",
			fmt.Sprint(user.ID),
			0,
			user.ID,
		); err != nil {
			panic(err) // TODO: Better error handling
		}
	}()

	return &api_v1.UpdateMeNotificationTimeSpanResponse{}, nil
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
