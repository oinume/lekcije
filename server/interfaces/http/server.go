package http

import (
	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/interfaces"
	"goji.io"
	"goji.io/pat"
	"gopkg.in/redis.v4"
)

type server struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewServer(args *interfaces.ServerArgs) *server {
	return &server{
		db:          args.DB,
		redisClient: args.RedisClient,
	}
}

func (s *server) Routes() *goji.Mux {
	mux := goji.NewMux()

	mux.HandleFunc(pat.Get("/static/*"), s.staticHandler())
	mux.HandleFunc(pat.Get("/"), s.indexHandler())
	mux.HandleFunc(pat.Get("/signup"), s.signupHandler())

	return mux
}
