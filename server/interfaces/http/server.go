package http

import (
	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/interfaces"
	"github.com/oinume/lekcije/server/interfaces/http/flash_message"
	"goji.io"
	"goji.io/pat"
	"gopkg.in/redis.v4"
)

type server struct {
	db                *gorm.DB
	flashMessageStore flash_message.Store
	redis             *redis.Client
}

func NewServer(args *interfaces.ServerArgs) *server {
	return &server{
		db:                args.DB,
		flashMessageStore: args.FlashMessageStore,
		redis:             args.Redis,
	}
}

func (s *server) Routes() *goji.Mux {
	mux := goji.NewMux()

	mux.HandleFunc(pat.Get("/static/*"), s.staticHandler())
	mux.HandleFunc(pat.Get("/"), s.indexHandler())
	mux.HandleFunc(pat.Get("/signup"), s.signupHandler())

	return mux
}
