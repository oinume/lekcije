package http

import (
	"net/http"

	"goji.io/v3"
	"goji.io/v3/pat"

	"github.com/oinume/lekcije/backend/usecase"
)

type SchedulerServer struct {
	notifier *usecase.Notifier
}

func NewSchedulerServer(notifier *usecase.Notifier) *SchedulerServer {
	return &SchedulerServer{
		notifier: notifier,
	}
}

func (s *SchedulerServer) Setup(mux *goji.Mux) {
	mux.HandleFunc(pat.Get("/scheduler/notifier"), s.notifierHandler)
}

func (s *SchedulerServer) notifierHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: fetch users
	s.notifier.SendNotification()
}
