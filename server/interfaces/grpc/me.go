package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/oinume/lekcije/server/logger"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/server/ga_measurement"

	"github.com/jinzhu/gorm"
	api_v1 "github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
	"github.com/oinume/lekcije/server/context_data"
	"github.com/oinume/lekcije/server/interfaces"
	"github.com/oinume/lekcije/server/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/redis.v4"
)

type apiV1Server struct {
	db                  *gorm.DB
	redis               *redis.Client
	gaMeasurementClient ga_measurement.Client
}

func RegisterAPIV1Server(server *grpc.Server, args *interfaces.ServerArgs) {
	api_v1.RegisterAPIServer(server, &apiV1Server{
		db:                  args.DB,
		redis:               args.Redis,
		gaMeasurementClient: args.GAMeasurementClient,
	})
}

func (s *apiV1Server) GetMe(
	ctx context.Context,
	in *api_v1.GetMeRequest,
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

func (s *apiV1Server) GetMeEmail(
	ctx context.Context, in *api_v1.GetMeEmailRequest,
) (*api_v1.GetMeEmailResponse, error) {
	user, err := authenticateFromContext(ctx, s.db)
	if err != nil {
		return nil, err
	}
	return &api_v1.GetMeEmailResponse{Email: user.Email}, nil
}

func (s *apiV1Server) UpdateMeEmail(
	ctx context.Context, in *api_v1.UpdateMeEmailRequest,
) (*api_v1.UpdateMeEmailResponse, error) {
	user, err := authenticateFromContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	// TODO: better validation
	email := in.Email
	if email == "" || !validateEmail(email) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Email")
	}

	userService := model.NewUserService(s.db)
	if err := userService.UpdateEmail(user, in.Email); err != nil {
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

func (s *apiV1Server) UpdateMeNotificationTimeSpan(
	ctx context.Context, in *api_v1.UpdateMeNotificationTimeSpanRequest,
) (*api_v1.UpdateMeNotificationTimeSpanResponse, error) {
	user, err := authenticateFromContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	// TODO: validation
	timeSpanService := model.NewNotificationTimeSpanService(s.db)
	timeSpans := timeSpanService.NewNotificationTimeSpansFromPB(user.ID, in.NotificationTimeSpans)
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

func (s *apiV1Server) sendGAMeasurementEvent(
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
		logger.App.Warn("SendEvent() failed", zap.Error(err))
	}
}

func authenticateFromContext(ctx context.Context, db *gorm.DB) (*model.User, error) {
	apiToken, err := context_data.GetAPIToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "No api token found")
	}
	userService := model.NewUserService(db)
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
