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

type MeService struct {
	appLogger                   *zap.Logger
	db                          *gorm.DB
	gaMeasurementUsecase        *usecase.GAMeasurement
	notificationTimeSpanUsecase *usecase.NotificationTimeSpan
	userUsecase                 *usecase.User
}

func NewMeService(
	db *gorm.DB, // TODO: Remove
	appLogger *zap.Logger,
	gaMeasurementUsecase *usecase.GAMeasurement,
	notificationTimeSpanUsecase *usecase.NotificationTimeSpan,
	userUsecase *usecase.User,
) api_v1.Me {
	return &MeService{
		appLogger:                   appLogger,
		db:                          db,
		gaMeasurementUsecase:        gaMeasurementUsecase,
		notificationTimeSpanUsecase: notificationTimeSpanUsecase,
		userUsecase:                 userUsecase,
	}
}

func (s *MeService) Ping(
	_ context.Context,
	_ *api_v1.PingRequest,
) (*api_v1.PingResponse, error) {
	return &api_v1.PingResponse{}, nil
}

func (s *MeService) GetMe(
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
	}, nil
}

func (s *MeService) GetEmail(
	ctx context.Context,
	_ *api_v1.GetEmailRequest,
) (*api_v1.GetEmailResponse, error) {
	user, err := authenticateFromContext(ctx, s.db)
	if err != nil {
		return nil, err
	}
	return &api_v1.GetEmailResponse{Email: user.Email}, nil
}

func (s *MeService) UpdateEmail(
	ctx context.Context,
	request *api_v1.UpdateEmailRequest,
) (*api_v1.UpdateEmailResponse, error) {
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

	return &api_v1.UpdateEmailResponse{}, nil
}

func (s *MeService) UpdateNotificationTimeSpan(
	ctx context.Context,
	request *api_v1.UpdateNotificationTimeSpanRequest,
) (*api_v1.UpdateNotificationTimeSpanResponse, error) {
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

	return &api_v1.UpdateNotificationTimeSpanResponse{}, nil
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
