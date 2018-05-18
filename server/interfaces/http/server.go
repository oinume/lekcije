package http

import (
	"goji.io"
	"goji.io/pat"
)

type server struct {
}

func NewServer() *server {
	return &server{}
}

func (s *server) Routes() *goji.Mux {
	mux := goji.NewMux()

	mux.HandleFunc(pat.Get("/static/*"), s.staticHandler())
	mux.HandleFunc(pat.Get("/"), s.indexHandler())
	mux.HandleFunc(pat.Get("/signup"), s.signupHandler())

	return mux
}
