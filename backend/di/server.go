package di

import (
	"database/sql"
	"net/http"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
	ihttp "github.com/oinume/lekcije/backend/interface/http"
	api_v1 "github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
)

func NewOAuthServer(
	appLogger *zap.Logger,
	db *sql.DB,
	gaMeasurementClient ga_measurement.Client,
	senderHTTPClient *http.Client,
) *ihttp.OAuthServer {
	return ihttp.NewOAuthServer(
		appLogger,
		gaMeasurementClient,
		NewGAMeasurementUsecase(gaMeasurementClient),
		senderHTTPClient,
		NewUserUsecase(db),
		NewUserAPITokenUsecase(db),
	)
}

func NewUserServer(
	appLogger *zap.Logger,
	db *gorm.DB,
	gaMeasurementClient ga_measurement.Client,
) api_v1.TwirpServer {
	userService := ihttp.NewUserService(
		db, appLogger,
		NewGAMeasurementUsecase(gaMeasurementClient),
		NewNotificationTimeSpanUsecase(db.DB()),
		NewUserUsecase(db.DB()),
	)
	return api_v1.NewUserServer(userService)
}
