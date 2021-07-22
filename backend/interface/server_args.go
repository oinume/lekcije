package interfaces

import (
	"database/sql"
	"net/http"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
	"github.com/oinume/lekcije/backend/interface/http/flash_message"
)

type ServerArgs struct {
	AccessLogger        *zap.Logger
	AppLogger           *zap.Logger
	DB                  *sql.DB
	FlashMessageStore   flash_message.Store
	GAMeasurementClient ga_measurement.Client
	GormDB              *gorm.DB
	SenderHTTPClient    *http.Client
}
