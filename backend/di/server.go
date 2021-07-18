package di

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/ga_measurement"
	ihttp "github.com/oinume/lekcije/backend/interface/http"
)

func NewOAuthServer(
	appLogger *zap.Logger,
	db *sql.DB,
	gaMeasurementClient ga_measurement.Client,
) *ihttp.OAuthServer {
	userUsecase := NewUserUsecase(db)
	userAPITokenUsecase := NewUserAPITokenUsecase(db)
	return ihttp.NewOAuthServer(appLogger, gaMeasurementClient, userUsecase, userAPITokenUsecase)
}
