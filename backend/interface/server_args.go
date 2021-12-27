package interfaces

import (
	"database/sql"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/rollbar/rollbar-go"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
)

type ServerArgs struct {
	AccessLogger        *zap.Logger
	AppLogger           *zap.Logger
	DB                  *sql.DB
	GAMeasurementClient ga_measurement.Client
	GormDB              *gorm.DB
	RollbarClient       *rollbar.Client
	SenderHTTPClient    *http.Client
}
