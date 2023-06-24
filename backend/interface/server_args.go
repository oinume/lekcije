package interfaces

import (
	"database/sql"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/jinzhu/gorm"
	"github.com/rollbar/rollbar-go"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
)

type ServerArgs struct {
	AccessLogger        *zap.Logger
	AppLogger           *zap.Logger
	DB                  *sql.DB
	FirebaseAuthClient  *auth.Client
	GAMeasurementClient ga_measurement.Client
	GormDB              *gorm.DB
	RollbarClient       *rollbar.Client
	SenderHTTPClient    *http.Client
}
