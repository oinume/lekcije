package http

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/ga_measurement"

	"github.com/jinzhu/gorm"
	"gopkg.in/redis.v4"

	"github.com/oinume/lekcije/backend/interfaces"
	"github.com/oinume/lekcije/backend/interfaces/http/flash_message"
)

type server struct {
	accessLogger        *zap.Logger
	appLogger           *zap.Logger
	db                  *gorm.DB
	flashMessageStore   flash_message.Store
	redis               *redis.Client
	senderHTTPClient    *http.Client
	gaMeasurementClient ga_measurement.Client
}

func NewServer(args *interfaces.ServerArgs) *server {
	return &server{
		accessLogger:        args.AccessLogger,
		appLogger:           args.AppLogger,
		db:                  args.DB,
		flashMessageStore:   args.FlashMessageStore,
		redis:               args.Redis,
		senderHTTPClient:    args.SenderHTTPClient,
		gaMeasurementClient: args.GAMeasurementClient,
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
