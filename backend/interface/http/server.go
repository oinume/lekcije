package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
	interfaces "github.com/oinume/lekcije/backend/interface"
	"github.com/oinume/lekcije/backend/usecase"
)

type server struct {
	accessLogger        *zap.Logger
	appLogger           *zap.Logger
	db                  *gorm.DB
	errorRecorder       *usecase.ErrorRecorder
	senderHTTPClient    *http.Client
	gaMeasurementClient ga_measurement.Client
	userAPITokenUsecase *usecase.UserAPIToken
}

func NewServer(
	args *interfaces.ServerArgs,
	errorRecorder *usecase.ErrorRecorder,
	userAPITokenUsecase *usecase.UserAPIToken,
) *server {
	return &server{
		accessLogger:        args.AccessLogger,
		appLogger:           args.AppLogger,
		db:                  args.GormDB,
		errorRecorder:       errorRecorder,
		senderHTTPClient:    args.SenderHTTPClient,
		gaMeasurementClient: args.GAMeasurementClient,
		userAPITokenUsecase: userAPITokenUsecase,
	}
}

func (s *server) sendGAMeasurementEvent(
	ctx context.Context,
	category,
	action,
	label string,
	value int64,
	userID uint32,
) {
	err := s.gaMeasurementClient.SendEvent(
		ctx,
		context_data.MustGAMeasurementEvent(ctx),
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
