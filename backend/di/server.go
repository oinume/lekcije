package di

import (
	"database/sql"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/rollbar/rollbar-go"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
	ihttp "github.com/oinume/lekcije/backend/interface/http"
	api_v1 "github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
)

func NewOAuthServer(
	appLogger *zap.Logger,
	db *sql.DB,
	gaMeasurementClient ga_measurement.Client,
	rollbarClient *rollbar.Client,
	senderHTTPClient *http.Client,
) *ihttp.OAuthServer {
	return ihttp.NewOAuthServer(
		appLogger,
		NewErrorRecorderUsecase(appLogger, rollbarClient),
		gaMeasurementClient,
		NewGAMeasurementUsecase(gaMeasurementClient),
		senderHTTPClient,
		NewUserUsecase(db),
		NewUserAPITokenUsecase(db),
	)
}

func NewMeServer(
	appLogger *zap.Logger,
	db *gorm.DB,
	errorRecorderHooks *twirp.ServerHooks,
	gaMeasurementClient ga_measurement.Client,
) api_v1.TwirpServer {
	meService := ihttp.NewMeService(
		db, appLogger,
		NewGAMeasurementUsecase(gaMeasurementClient),
		NewNotificationTimeSpanUsecase(db.DB()),
		NewUserUsecase(db.DB()),
	)
	return api_v1.NewMeServer(
		meService,
		twirp.WithServerHooks(errorRecorderHooks),
	)
}
