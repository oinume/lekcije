package registry

import (
	"database/sql"
	"net/http"

	"github.com/rollbar/rollbar-go"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
	ihttp "github.com/oinume/lekcije/backend/interface/http"
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
