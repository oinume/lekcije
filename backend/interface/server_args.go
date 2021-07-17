package interfaces

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/ga_measurement"
	"github.com/oinume/lekcije/backend/interface/http/flash_message"
)

type ServerArgs struct {
	AccessLogger        *zap.Logger
	AppLogger           *zap.Logger
	DB                  *gorm.DB
	FlashMessageStore   flash_message.Store
	SenderHTTPClient    *http.Client
	GAMeasurementClient ga_measurement.Client
}
