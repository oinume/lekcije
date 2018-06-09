package http

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/interfaces"
	"github.com/oinume/lekcije/server/interfaces/http/flash_message"
	"gopkg.in/redis.v4"
)

type server struct {
	db                *gorm.DB
	flashMessageStore flash_message.Store
	redis             *redis.Client
	senderHTTPClient  *http.Client
}

func NewServer(args *interfaces.ServerArgs) *server {
	return &server{
		db:                args.DB,
		flashMessageStore: args.FlashMessageStore,
		redis:             args.Redis,
		senderHTTPClient:  args.SenderHTTPClient,
	}
}
