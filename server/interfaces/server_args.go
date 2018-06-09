package interfaces

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/interfaces/http/flash_message"
	"gopkg.in/redis.v4"
)

type ServerArgs struct {
	DB                *gorm.DB
	FlashMessageStore flash_message.Store
	Redis             *redis.Client
	SenderHTTPClient  *http.Client
}
